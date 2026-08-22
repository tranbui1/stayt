import React from "react";
import { View } from "react-native";
import { colors } from "../theme/colors";

export default function Dot({ size = 14, color = colors.dot, style }) {
  return (
    <View
      style={[
        { width: size, height: size, borderRadius: size / 2, backgroundColor: color },
        style,
      ]}
    />
  );
}
