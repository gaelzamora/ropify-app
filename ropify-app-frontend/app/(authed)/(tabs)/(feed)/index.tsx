import { Text, View } from "react-native";
import React from "react";

export default function FeedScreen() {
    return (
        <View style={{
            flex: 1,
            alignItems: "center",
            justifyContent: "center",
        }}>
            <Text style={{
                textAlign: "center",
                fontSize: 20,
                color: "#e85a5a",
                textTransform: "uppercase",
                width: 300,
                fontWeight: "800"
            }}>Coming soon</Text>
        </View>
    )
}