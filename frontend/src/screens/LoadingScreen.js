import React from "react";
import { View, Text, StyleSheet } from "react-native";
import Particle from "../components/Particle";
import Star from "../components/Star";
import Dot from "../components/Dot";
import { colors } from "../theme/colors";

const PARTICLES = [
  { type: "star", top: "10%", left: "80%", size: 30, delay: 0 },
  { type: "dot", top: "20%", left: "12%", size: 14, delay: 250 },
  { type: "star", top: "36%", left: "13%", size: 26, delay: 500 },
  { type: "dot", top: "58%", left: "10%", size: 14, delay: 750 },
  { type: "star", top: "64%", left: "83%", size: 28, delay: 150 },
  { type: "star", top: "73%", left: "34%", size: 30, delay: 600 },
  { type: "dot", top: "84%", left: "80%", size: 14, delay: 400 },
];

export default function LoadingScreen() {
  return (
    <View style={styles.container}>
      {PARTICLES.map((p, i) => (
        <Particle
          key={i}
          delay={p.delay}
          duration={1000}
          bounceHeight={10}
          style={[styles.particle, { top: p.top, left: p.left }]}
        >
          {p.type === "star" ? (
            <Star size={p.size} />
          ) : (
            <Dot size={p.size} />
          )}
        </Particle>
      ))}

      <View style={styles.textBlock}>
        <View style={styles.badge}>
          <Text style={styles.badgeText}>+1</Text>
        </View>
        <Text style={styles.title}>Someone Thought{"\n"}About You Today</Text>
      </View>
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
  particle: {
    position: "absolute",
  },
  textBlock: {
    alignItems: "center",
  },
  badge: {
    position: "absolute",
    top: -14,
    right: -28,
    backgroundColor: colors.badge,
    borderRadius: 10,
    paddingHorizontal: 8,
    paddingVertical: 2,
  },
  badgeText: {
    color: "#ffffff",
    fontSize: 12,
    fontWeight: "600",
  },
  title: {
    color: colors.text,
    fontSize: 16,
    textAlign: "center",
    lineHeight: 22,
  },
});
