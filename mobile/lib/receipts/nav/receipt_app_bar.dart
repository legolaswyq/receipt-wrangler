import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';

import '../../utils/forms.dart';
import '../../utils/group.dart';
import '../../utils/receipts.dart';

class ReceiptAppBar extends StatefulWidget implements PreferredSizeWidget {
  const ReceiptAppBar({super.key, this.actions});

  final List<Widget>? actions;

  @override
  State<ReceiptAppBar> createState() => _ReceiptAppBar();

  @override
  Size get preferredSize => AppBar().preferredSize;
}

class _ReceiptAppBar extends State<ReceiptAppBar> {
  // Pinned to the app's single persistent group context rather than
  // receiptModel.receipt.groupId — a receipt found via search can belong to a
  // different real group, and following that group's id here would silently
  // land the user on that group's (possibly widget-lighter) routes on the next
  // bottom-nav tap. See utils/group.dart `defaultGroupId`.
  String buildBackUrl(BuildContext context, WranglerFormState formState) {
    if (formState == WranglerFormState.add) {
      return "/groups";
    } else {
      return "/groups/${defaultGroupId(context)}/receipts";
    }
  }

  @override
  Widget build(BuildContext context) {
    var uri = GoRouter.of(context).routeInformationProvider.value.uri;
    var formState = getFormState(uri.toString());
    List<Widget> actions = [...widget.actions ?? []];

    // TODO: fix add redirect
    return Consumer<ReceiptModel>(
        builder: (context, receiptModel, child) => TopAppBar(
              titleText: getTitleText(formState, receiptModel.receipt.name),
              leadingArrowRedirect: buildBackUrl(context, formState),
              leadingArrowPop: false,
              onLeadingArrowPressed: () {
                receiptModel.resetModel();
              },
              actions: actions,
              hideAvatar: true,
            ));
  }
}
