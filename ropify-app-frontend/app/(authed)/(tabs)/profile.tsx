import { useAuth } from "@/context/AuthContext";
import React from "react";
import { TouchableOpacity, View, Text } from "react-native";

export default function ProfileScreen() {
    const { logout } = useAuth()

    return (
        <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}>
            <TouchableOpacity
                style={{ 
                    backgroundColor: "#222", 
                    paddingVertical: 15, 
                    paddingHorizontal: 35,
                    borderRadius: 10
                }}
                onPress={logout}
            >
                <Text style={{ color: "white", fontWeight: "700" }}>Logout</Text>
            </TouchableOpacity>
        </View>
    )
}