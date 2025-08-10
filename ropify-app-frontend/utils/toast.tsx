import { StyleSheet, View } from "react-native";
import Toast, { BaseToast, ErrorToast } from "react-native-toast-message";

// Custom config
export const toastConfig = {
  success: (props: any) => (
    <View style={styles.container}>
      <BaseToast
        {...props}
        style={{ borderLeftColor: "#222" }} 
        contentContainerStyle={{ paddingHorizontal: 15 }}
        text1Style={{
          fontSize: 16,
          fontWeight: "bold"
        }}
        text2Style={{
          fontSize: 14
        }}
      />
    </View>
  ),
  error: (props: any) => (
    <View style={styles.container} >
      <ErrorToast
        {...props}
        style={{ borderLeftColor: "#ab114c" }} 
        text1Style={{
          fontSize: 16,
          fontWeight: "bold"
        }}
        text2Style={{
          fontSize: 14
        }}
      />
    </View>
  ),
};

const styles = StyleSheet.create({
  container: {
    zIndex: 9999,
  }
})