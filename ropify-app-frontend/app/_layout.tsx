import { Slot } from "expo-router";
import { StatusBar } from "react-native";

export default function Root() {
  return (
    <>
      <StatusBar />
      <Slot />
    </>
  )
}