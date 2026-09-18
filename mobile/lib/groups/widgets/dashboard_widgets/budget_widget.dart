import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/utils/currency.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../../client/client.dart';
import '../constants/text_styles.dart';

/// Monthly Budget widget — the mobile counterpart of the desktop Budget
/// widget. Shows the current month's income/spent/net summary plus a
/// per-category progress bar with tap-to-edit targets. Kept in sync with
/// desktop/src/dashboard/budget/budget.component.*.
class BudgetWidget extends StatefulWidget {
  const BudgetWidget({super.key, required this.dashboardWidget});

  final api.Widget dashboardWidget;

  @override
  State<BudgetWidget> createState() => _BudgetWidgetState();
}

class _BudgetWidgetState extends State<BudgetWidget> {
  late Future<api.BudgetData?> _future;
  bool _isInitialized = false;

  int? _editingCategoryId;
  final _editController = TextEditingController();

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (!_isInitialized) {
      _loadData();
      _isInitialized = true;
    }
  }

  @override
  void dispose() {
    _editController.dispose();
    super.dispose();
  }

  int get _groupId => int.tryParse(getGroupId(context)) ?? 0;

  void _loadData() {
    _future = OpenApiClient.client
        .getBudgetApi()
        .getBudgetData(groupId: _groupId)
        .then((response) => response.data);
  }

  double _progressPercent(api.BudgetCategory category) {
    final target = category.target ?? 0;
    if (target <= 0) return 0;
    final spent = category.spent ?? 0;
    return (spent / target).clamp(0, 1);
  }

  bool _isOver(api.BudgetCategory category) => category.over ?? false;

  void _startEdit(api.BudgetCategory category) {
    setState(() {
      _editingCategoryId = category.categoryId;
      _editController.text = (category.target ?? 0).toStringAsFixed(2);
    });
  }

  Future<void> _saveTarget(api.BudgetCategory category) async {
    final categoryId = category.categoryId;
    final amount = double.tryParse(_editController.text);
    setState(() => _editingCategoryId = null);

    if (categoryId == null || amount == null || amount <= 0) {
      return;
    }

    final command = api.UpsertCategoryBudgetCommand((b) => b
      ..categoryId = categoryId
      ..amount = amount.toStringAsFixed(2));

    await OpenApiClient.client
        .getBudgetApi()
        .upsertBudget(groupId: _groupId, upsertCategoryBudgetCommand: command);

    if (!mounted) return;
    setState(_loadData);
  }

  Widget _buildSummaryItem(String label, String value, {bool negative = false}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
        Text(
          value,
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: negative ? Colors.red[700] : null,
          ),
        ),
      ],
    );
  }

  Widget _buildCategoryRow(api.BudgetCategory category) {
    final isOver = _isOver(category);
    final isEditing = _editingCategoryId == category.categoryId;
    final overColor = Colors.red[700];

    return Padding(
      key: Key('budget-row-${category.categoryId}'),
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Expanded(
                child: Text(
                  category.name ?? '',
                  style: const TextStyle(fontWeight: FontWeight.w500),
                ),
              ),
              if (isEditing)
                SizedBox(
                  width: 90,
                  child: TextField(
                    controller: _editController,
                    autofocus: true,
                    keyboardType:
                        const TextInputType.numberWithOptions(decimal: true),
                    textAlign: TextAlign.right,
                    decoration: const InputDecoration(isDense: true),
                    onSubmitted: (_) => _saveTarget(category),
                    onTapOutside: (_) => _saveTarget(category),
                  ),
                )
              else
                GestureDetector(
                  onTap: () => _startEdit(category),
                  child: Text(
                    '${formatCurrency(context, (category.spent ?? 0).toString())} / '
                    '${formatCurrency(context, (category.target ?? 0).toString())}',
                    style: TextStyle(
                      color: isOver ? overColor : null,
                      fontWeight: isOver ? FontWeight.w600 : FontWeight.normal,
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 4),
          ClipRRect(
            borderRadius: BorderRadius.circular(4),
            child: LinearProgressIndicator(
              value: _progressPercent(category),
              minHeight: 8,
              backgroundColor: Colors.grey[200],
              color: isOver ? overColor : Theme.of(context).colorScheme.primary,
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<api.BudgetData?>(
      future: _future,
      builder: (context, snapshot) {
        final isLoading = snapshot.connectionState != ConnectionState.done;
        final data = snapshot.data;

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  widget.dashboardWidget.name ?? 'Budget',
                  style: dashboardWidgetNameStyle,
                ),
                if (data?.month != null)
                  Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.grey[200],
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(data!.month!, style: const TextStyle(fontSize: 12)),
                  ),
              ],
            ),
            const SizedBox(height: 10),
            if (isLoading)
              const Padding(
                padding: EdgeInsets.symmetric(vertical: 16),
                child: Center(child: CircularProgressIndicator()),
              )
            else if (data == null)
              Text('Failed to load budget data',
                  style: TextStyle(color: Colors.grey[600]))
            else ...[
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  _buildSummaryItem('Income',
                      formatCurrency(context, (data.income ?? 0).toString()) ?? ''),
                  _buildSummaryItem('Spent',
                      formatCurrency(context, (data.spent ?? 0).toString()) ?? ''),
                  _buildSummaryItem(
                    'Net',
                    formatCurrency(context, (data.net ?? 0).toString()) ?? '',
                    negative: (data.net ?? 0) < 0,
                  ),
                ],
              ),
              const SizedBox(height: 16),
              if (data.categories.isEmpty)
                Text('No category budgets set yet.',
                    style: TextStyle(color: Colors.grey[600]))
              else
                ...data.categories.map(_buildCategoryRow),
              if ((data.untracked ?? 0) > 0) ...[
                const Divider(),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text('Untracked', style: TextStyle(color: Colors.grey[600])),
                    Text(
                      formatCurrency(context, (data.untracked ?? 0).toString()) ??
                          '',
                      style: TextStyle(color: Colors.grey[600]),
                    ),
                  ],
                ),
              ],
            ],
          ],
        );
      },
    );
  }
}
