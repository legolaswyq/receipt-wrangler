import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import '../models/group_model.dart';

String getGroupId(BuildContext context) {
  var extraMap =
      (GoRouterState.of(context).extra ?? {}) as Map<dynamic, dynamic>;
  return GoRouterState.of(context).pathParameters["groupId"] ??
      extraMap["groupId"] ??
      "0";
}

/// The app's single, persistent group context (the group concept is hidden
/// from the UI — see `singleGroupRedirect`), as opposed to [getGroupId]'s
/// current-route group, which can differ: viewing a receipt that belongs to
/// a different real group and tapping back lands you on THAT group's routes
/// (`ReceiptAppBar.buildBackUrl` uses the receipt's own `groupId`, correctly,
/// since the receipt itself is scoped there). Persistent nav actions (the
/// group bottom nav's Dashboards/Receipts destinations) must not follow that
/// drift — otherwise tapping "Dashboards" from a real group's receipt list
/// silently shows THAT group's dashboard, which may be missing widgets the
/// default group's dashboard has (observed: a real group's dashboard lacked
/// the Category Breakdown widget the aggregate "All" group's dashboard has).
/// Falls back to [getGroupId] only if no group has loaded yet, so navigation
/// still works during the narrow window before `GroupModel` is populated.
String defaultGroupId(BuildContext context) {
  final defaultGroup =
      Provider.of<GroupModel>(context, listen: false).defaultGroup;
  return defaultGroup?.id.toString() ?? getGroupId(context);
}

String? getGroupByIdWithRouter(GoRouter router) {
  return router.routerDelegate.currentConfiguration.pathParameters["groupId"];
}
