import { AuthenticationProvider } from "@/context/AuthContext";
import { Slot } from "expo-router";
import { StatusBar, View } from "react-native";
import Toast from "react-native-toast-message"
import { StyleSheet } from "react-native";
import { toastConfig } from "@/utils/toast";

export default function Root() {
  return (
    <>
      <StatusBar />
      <View style={styles.container}>
        <Toast config={toastConfig} />
      </View>
      <AuthenticationProvider>
        <Slot />
      </AuthenticationProvider>
    </>
  )
}

const styles = StyleSheet.create({
  container: {
    zIndex: 10000
  },
});