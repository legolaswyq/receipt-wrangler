import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/nav/group/group_app_bar.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';

import '../helpers/auth_test_helpers.dart';

api.Claims _claims(String displayName) =>
    api.Claims((b) => b..displayName = displayName);

// The group concept is hidden, so the app bar titles by section (Dashboard /
// Receipts) rather than by group name, and shows no back arrow (there is no
// group-select screen to return to). GroupAppBar no longer reads GroupModel;
// the provider is still supplied because TopAppBar's avatar/menu reads other
// models from the same tree.
Future<void> _pumpAt(
  WidgetTester tester, {
  required String initialLocation,
}) async {
  final groupModel = MockGroupModel();

  final authModel = MockAuthModel();
  when(() => authModel.claims).thenReturn(_claims('Admin'));

  final router = GoRouter(
    initialLocation: initialLocation,
    routes: [
      GoRoute(
        path: '/groups/:groupId/dashboards',
        builder: (_, __) => const Scaffold(appBar: GroupAppBar()),
      ),
      GoRoute(
        path: '/groups/:groupId/receipts',
        builder: (_, __) => const Scaffold(appBar: GroupAppBar()),
      ),
    ],
  );

  await tester.pumpWidget(MultiProvider(
    providers: [
      Provider<AuthModel>.value(value: authModel),
      ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
      Provider<GroupModel>.value(value: groupModel),
      ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
    ],
    child: MaterialApp.router(routerConfig: router),
  ));
  await tester.pump();
}

Finder _titleInside(String text) => find.descendant(
      of: find.byType(GroupAppBar),
      matching: find.text(text),
    );

void main() {
  setUpAll(() {
    registerFallbackValue('');
    Provider.debugCheckInvalidValueType = null;
  });

  testWidgets('dashboards route → title "Dashboard"', (tester) async {
    await _pumpAt(tester, initialLocation: '/groups/1/dashboards');
    expect(_titleInside('Dashboard'), findsOneWidget);
    expect(tester.takeException(), isNull);
  });

  testWidgets('receipts route → title "Receipts"', (tester) async {
    await _pumpAt(tester, initialLocation: '/groups/1/receipts');
    expect(_titleInside('Receipts'), findsOneWidget);
    expect(tester.takeException(), isNull);
  });

  testWidgets('no back arrow is shown (group-select is gone)', (tester) async {
    await _pumpAt(tester, initialLocation: '/groups/1/dashboards');
    expect(find.byIcon(Icons.arrow_back), findsNothing);
  });
}
