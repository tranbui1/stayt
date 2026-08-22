import React, { useEffect, useState } from "react";
import { StatusBar } from "expo-status-bar";
import LogoScreen from "./src/screens/LogoScreen";
import LoadingScreen from "./src/screens/LoadingScreen";

const LOGO_DURATION_MS = 1500;

export default function App() {
  const [screen, setScreen] = useState("logo");

  useEffect(() => {
    const timer = setTimeout(() => setScreen("loading"), LOGO_DURATION_MS);
    return () => clearTimeout(timer);
  }, []);

  return (
    <>
      <StatusBar style="dark" />
      {screen === "logo" ? <LogoScreen /> : <LoadingScreen />}
    </>
  );
}
