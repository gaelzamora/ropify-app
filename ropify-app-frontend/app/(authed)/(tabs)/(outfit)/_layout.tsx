import { Stack } from "expo-router";
import React from "react";

export default function OutfitLayout() {
    return (
        <Stack screenOptions={{ headerShown: false }}>
            <Stack.Screen name="index" />
        </Stack>
    )
}