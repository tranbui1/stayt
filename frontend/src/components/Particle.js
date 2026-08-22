import React, { useEffect, useRef } from "react";
import { Animated, Easing } from "react-native";

export default function Particle({
  children,
  style,
  delay = 0,
  duration = 1000,
  bounceHeight = 10,
}) {
  const translateY = useRef(new Animated.Value(0)).current;
  const opacity = useRef(new Animated.Value(0.25)).current;

  useEffect(() => {
    const cycle = Animated.parallel([
      Animated.sequence([
        Animated.timing(translateY, {
          toValue: -bounceHeight,
          duration: duration / 2,
          easing: Easing.out(Easing.quad),
          useNativeDriver: true,
        }),
        Animated.timing(translateY, {
          toValue: 0,
          duration: duration / 2,
          easing: Easing.in(Easing.quad),
          useNativeDriver: true,
        }),
      ]),
      Animated.sequence([
        Animated.timing(opacity, {
          toValue: 1,
          duration: duration / 2,
          useNativeDriver: true,
        }),
        Animated.timing(opacity, {
          toValue: 0.25,
          duration: duration / 2,
          useNativeDriver: true,
        }),
      ]),
    ]);

    const loop = Animated.loop(Animated.sequence([Animated.delay(delay), cycle]));
    loop.start();
    return () => loop.stop();
  }, [bounceHeight, delay, duration, opacity, translateY]);

  return (
    <Animated.View style={[style, { opacity, transform: [{ translateY }] }]}>
      {children}
    </Animated.View>
  );
}
