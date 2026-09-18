import 'package:flutter/material.dart';

/// Parses a `#RRGGBB` (or bare `RRGGBB`) hex string into an opaque [Color].
/// Returns null for null/empty/malformed input so callers can fall back to a
/// palette (or omit the swatch). Shared by the dashboard pie chart and the
/// spending table so both resolve the API's stored category colors the same way.
Color? hexToColor(String? hex) {
  if (hex == null) return null;
  var value = hex.trim();
  if (value.startsWith('#')) value = value.substring(1);
  if (value.length != 6) return null;
  final rgb = int.tryParse(value, radix: 16);
  if (rgb == null) return null;
  return Color(0xFF000000 | rgb);
}
