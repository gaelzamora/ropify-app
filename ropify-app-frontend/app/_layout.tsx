import { AuthenticationProvider } from "@/context/AuthContext";
import { Slot } from "expo-router";
import { StatusBar, View } from "react-native";

export default function Root() {
  return (
    <>
      <StatusBar />
      <AuthenticationProvider>
        <Slot />
      </AuthenticationProvider>
    </>
  )
}