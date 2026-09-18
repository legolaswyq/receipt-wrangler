import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';

class GroupAppBar extends StatefulWidget implements PreferredSizeWidget {
  const GroupAppBar({super.key});

  @override
  State<GroupAppBar> createState() => _GroupAppBar();

  @override
  Size get preferredSize => AppBar().preferredSize;
}

class _GroupAppBar extends State<GroupAppBar> {
  /// The group concept is hidden, so the app bar reads by section rather than
  /// by group name — mirroring the desktop's Dashboard / Receipts tabs. No back
  /// arrow, because there is no group-select screen to return to.
  String _titleForRoute(BuildContext context) {
    final uri =
        GoRouter.of(context).routeInformationProvider.value.uri.toString();
    if (uri.contains("dashboards")) {
      return "Dashboard";
    }
    return "Receipts";
  }

  @override
  Widget build(BuildContext context) {
    return TopAppBar(titleText: _titleForRoute(context));
  }
}
