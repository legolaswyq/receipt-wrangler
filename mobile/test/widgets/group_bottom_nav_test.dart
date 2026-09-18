import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/nav/group/group_bottom_nav.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

import '../helpers/auth_test_helpers.dart';
import '../helpers/permission_test_helpers.dart';

/// The group-context bottom nav hides its trailing Search destination when the
/// caller lacks `app.receipts.search`. A real `GoRouter` ancestor is required
/// because `setIndexSelected` reads `GoRouter.of(context)` during build.
void main() {
  setUpAll(() {
    registerFallbackValue('');
    Provider.debugCheckInvalidValueType = null;
  });

  Widget harness(PermissionsModel permissions) {
    final router = GoRouter(
      initialLocation: '/groups/1/receipts',
      routes: [
        GoRoute(
          path: '/groups/:groupId/receipts',
          builder: (context, state) =>
              const Scaffold(bottomNavigationBar: GroupBottomNav()),
        ),
      ],
    );
    return ChangeNotifierProvider<PermissionsModel>.value(
      value: permissions,
      child: MaterialApp.router(routerConfig: router),
    );
  }

  testWidgets('hides the Search destination without app.receipts.search',
      (tester) async {
    await tester.pumpWidget(harness(seededPermissions()));
    await tester.pump();

    expect(find.byIcon(Icons.search), findsNothing);
    expect(find.text('Search'), findsNothing);
    // The other destinations remain.
    expect(find.byIcon(Icons.dashboard), findsOneWidget);
    expect(find.byIcon(Icons.add), findsOneWidget);
    expect(find.byIcon(Icons.receipt), findsOneWidget);
  });

  testWidgets('shows the Search destination with app.receipts.search',
      (tester) async {
    await tester.pumpWidget(harness(
      seededPermissions(app: [Permission.appPeriodReceiptsPeriodSearch]),
    ));
    await tester.pump();

    expect(find.byIcon(Icons.search), findsOneWidget);
    expect(find.text('Search'), findsOneWidget);
  });

  // Regression coverage for a bug found in manual testing: viewing a receipt
  // that belongs to a different real group (via search) and tapping back
  // lands you on THAT group's routes (ReceiptAppBar.buildBackUrl uses the
  // receipt's own groupId, correctly). If the bottom nav's Dashboards/Receipts
  // destinations then just read the current route's groupId, tapping them
  // silently shows that OTHER group's dashboard/receipts -- which may be
  // missing widgets the default group's dashboard has (observed: the Category
  // Breakdown widget). The destinations must stay pinned to the app's single
  // persistent default group instead.
  testWidgets(
      'Dashboards destination navigates to the default group, not the '
      'current route\'s group', (tester) async {
    final defaultGroup = MockGroup();
    when(() => defaultGroup.id).thenReturn(2);
    final groupModel = MockGroupModel();
    when(() => groupModel.defaultGroup).thenReturn(defaultGroup);

    final router = GoRouter(
      initialLocation: '/groups/1/receipts',
      routes: [
        GoRoute(
          path: '/groups/:groupId/receipts',
          builder: (context, state) =>
              const Scaffold(bottomNavigationBar: GroupBottomNav()),
        ),
        GoRoute(
          path: '/groups/:groupId/dashboards',
          builder: (context, state) => Scaffold(
            body: Text('dashboards for ${state.pathParameters["groupId"]}'),
            bottomNavigationBar: const GroupBottomNav(),
          ),
        ),
      ],
    );

    await tester.pumpWidget(MultiProvider(
      providers: [
        ChangeNotifierProvider<PermissionsModel>.value(
          value: seededPermissions(),
        ),
        Provider<GroupModel>.value(value: groupModel),
      ],
      child: MaterialApp.router(routerConfig: router),
    ));
    await tester.pump();

    await tester.tap(find.byIcon(Icons.dashboard));
    await tester.pumpAndSettle();

    // Landed on the default group (2), not the route it was tapped from (1).
    expect(find.text('dashboards for 2'), findsOneWidget);
  });
}
