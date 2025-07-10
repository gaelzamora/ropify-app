import { Text, View } from "react-native";
import React from "react";
import { Ionicons } from "@expo/vector-icons";

export default function FeedScreen() {
    return (
        <View style={{
            flex: 1,
            alignItems: "center",
            justifyContent: "center",
            width: 200,
            marginHorizontal: "auto"
        }}>
            <Ionicons 
                name="rocket-outline"
                color={"#222"}
                size={70}
            />
            <Text style={{
                textAlign: "center",
                fontSize: 20,
                color: "#222",
                textTransform: "uppercase",
                width: 300,
                fontWeight: "800",
                marginTop: 10
            }}>Coming soon</Text>

            <Text
                style={{
                    textAlign: "center",
                    fontSize: 12,
                    color: "#888",
                    fontWeight: "800",
                    marginTop: 5
                }}
            >
                We’re working on something amazing. Stay tuned!
            </Text>
        </View>
    )
}