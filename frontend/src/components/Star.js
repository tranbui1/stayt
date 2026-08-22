import React from "react";
import { Text } from "react-native";
import { colors } from "../theme/colors";

export default function Star({ size = 40, color = colors.star, style }) {
  return (
    <Text style={[{ fontSize: size, lineHeight: size, color }, style]}>★</Text>
  );
}
