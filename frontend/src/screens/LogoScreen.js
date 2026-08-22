import React from "react";
import { View, StyleSheet } from "react-native";
import Star from "../components/Star";
import { colors } from "../theme/colors";

export default function LogoScreen() {
  return (
    <View style={styles.container}>
      <Star size={140} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    alignItems: "center",
    justifyContent: "center",
  },
});
