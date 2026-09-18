import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';

class UserAvatar extends StatefulWidget {
  const UserAvatar({super.key, this.userId});

  final String? userId;

  @override
  State<UserAvatar> createState() => _UserAvatar();
}

class _UserAvatar extends State<UserAvatar> {
  String _getAvatarText() {
    var authModel = Provider.of<AuthModel>(context, listen: true);
    var userModel = Provider.of<UserModel>(context, listen: true);

    // Degrade to a placeholder rather than crashing when the name isn't
    // available yet — e.g. a screen with the top app bar renders during a
    // cold start before claims/users have loaded, or a partial-failure
    // startup leaves claims null. Previously this force-unwrapped and threw
    // "Null check operator used on a null value".
    String? name;
    if (widget.userId?.isNotEmpty == true) {
      name = userModel.getUserById(widget.userId!)?.displayName;
    } else {
      name = authModel.claims?.displayName;
    }

    return (name != null && name.isNotEmpty) ? name[0].toUpperCase() : "?";
  }

  @override
  Widget build(BuildContext context) {
    return CircleAvatar(
      child: Text(_getAvatarText()),
    );
  }
}
