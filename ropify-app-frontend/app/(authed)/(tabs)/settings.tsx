import { useAuth } from "@/context/AuthContext";
import React from "react";
import { TouchableOpacity, View, Text } from "react-native";

export default function SettingsScreen() {
    const { logout } = useAuth()

    return (
        <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}>
            <TouchableOpacity
                style={{ backgroundColor: "#888", paddingVertical: 5, paddingHorizontal: 15 }}
                onPress={logout}
            >
                <Text>Logout</Text>
            </TouchableOpacity>
        </View>
    )
}