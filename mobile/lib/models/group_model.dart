import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class GroupModel extends ChangeNotifier {
  List<Group> _groups = [];

  List<Group> get groups => _groups;

  List<Group> get groupsWithoutAllGroup =>
      _groups.where((group) => !group.isAllGroup).toList();

  /// The group the app silently uses now that the group concept is hidden from
  /// the UI. Prefers the aggregate "All" group (a unified, group-free overview)
  /// and falls back to the first active group. Returns null only when the user
  /// belongs to no active group.
  Group? get defaultGroup {
    final active =
        _groups.where((group) => group.status == GroupStatus.ACTIVE).toList();
    if (active.isEmpty) {
      return null;
    }
    return active.firstWhere(
      (group) => group.isAllGroup,
      orElse: () => active.first,
    );
  }

  void setGroups(List<Group> groups) {
    _groups = groups;

    notifyListeners();
  }

  Group? getGroupById(String id) {
    try {
      return _groups.firstWhere((group) => group.id == int.tryParse(id));
    } catch (e) {
      return null;
    }
  }

  GroupReceiptSettings? getGroupReceiptSettings(int groupId) {
    if (groupId == 0) {
      return null;
    }

    return getGroupById(groupId.toString())?.groupReceiptSettings;
  }
}
