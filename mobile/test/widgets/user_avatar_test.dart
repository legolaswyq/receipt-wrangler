import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/user_avatar.dart';

import '../helpers/auth_test_helpers.dart';

Future<void> _pump(WidgetTester tester, AuthModel authModel) async {
  await tester.pumpWidget(MultiProvider(
    providers: [
      Provider<AuthModel>.value(value: authModel),
      Provider<UserModel>.value(value: MockUserModel()),
    ],
    child: const MaterialApp(home: Scaffold(body: UserAvatar())),
  ));
  await tester.pump();
}

void main() {
  setUpAll(() => Provider.debugCheckInvalidValueType = null);

  testWidgets('renders the first initial from the caller claims',
      (tester) async {
    final authModel = MockAuthModel();
    when(() => authModel.claims)
        .thenReturn(api.Claims((b) => b..displayName = 'Admin'));
    await _pump(tester, authModel);
    expect(find.text('A'), findsOneWidget);
    expect(tester.takeException(), isNull);
  });

  testWidgets('degrades to "?" without crashing when claims are null',
      (tester) async {
    final authModel = MockAuthModel();
    when(() => authModel.claims).thenReturn(null);
    await _pump(tester, authModel);
    expect(find.text('?'), findsOneWidget);
    expect(tester.takeException(), isNull);
  });
}
