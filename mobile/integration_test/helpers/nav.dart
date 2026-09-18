import 'package:flutter_test/flutter_test.dart';

import 'pump.dart';

/// Confirms the group context shell is ready. The group concept is hidden, so
/// login now lands directly inside the default group's shell (see
/// `singleGroupRedirect`) instead of on a group-select screen — there is no
/// group card to tap. [groupName] is retained so existing call sites read the
/// same, but is unused; a freshly-provisioned permission user belongs to a
/// single group, which is what the app lands on.
///
/// Waits for *hittability*, not bare existence: on iOS the Cupertino page
/// transitions (~400ms) keep the destination sliding after its widgets mount.
/// See mobile/CLAUDE.md "Three tap-flake patterns".
Future<void> enterGroup(WidgetTester tester, String groupName) async {
  // The group bottom nav ("Dashboards"/"Add"/"Receipts"/"Search") only mounts
  // inside the group context shell -- a stronger signal than a URL check.
  await pumpUntilFound(tester, find.text('Receipts').hitTestable());
  await _drain(tester);
}

/// Enters [groupName] and opens its Receipts tab, returning once [receiptName]
/// is on screen.
Future<void> openGroupReceipts(
  WidgetTester tester,
  String groupName,
  String receiptName,
) async {
  await enterGroup(tester, groupName);
  await _drain(tester);
  await tester.tap(find.text('Receipts').hitTestable());
  await pumpUntilFound(tester, find.text(receiptName));
}

/// Drains the tail of a page/sheet transition so a follow-up tap computes its
/// center from settled geometry.
Future<void> _drain(WidgetTester tester) async {
  for (int i = 0; i < 5; i++) {
    await tester.pump(const Duration(milliseconds: 100));
  }
}
