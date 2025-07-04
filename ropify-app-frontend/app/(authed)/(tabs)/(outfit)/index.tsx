import { Text, View } from "react-native";
import React from "react";
import { Ionicons } from "@expo/vector-icons";

export default function OutfitScreen() {
    return (
        <View style={{
            flex: 1,
            alignItems: "center",
            justifyContent: "center",
        }}>
            <Ionicons 
                name="rocket-outline"
                color={"#d85858"}
                size={70}
            />
            <Text style={{
                textAlign: "center",
                fontSize: 20,
                color: "#d85858",
                textTransform: "uppercase",
                width: 300,
                fontWeight: "800",
                marginTop: 10
            }}>Coming soon</Text>

            <Text
                style={{
                    textAlign: "center",
                    fontSize: 12,
                    color: "#d85858",
                    fontWeight: "800",
                    marginTop: 5
                }}
            >
                We’re working on something amazing. Stay tuned!
            </Text>
        </View>
    )
}