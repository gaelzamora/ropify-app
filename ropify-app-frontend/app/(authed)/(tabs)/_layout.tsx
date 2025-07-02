import { Ionicons } from "@expo/vector-icons";
import { Tabs } from "expo-router";
import React, { ComponentProps } from "react";
import { View, Text } from "react-native";

export default function TabLayout() {

    const tabs = [
        {
            name: "(feed)",
            displayName: "Home",
            icon: "home-outline",
        },
        {
            name: "(closet)",
            displayName: "Closet",
            icon: "folder-outline",
        },
        {
            name: "(outfit)",
            displayName: "Outfit",
            icon: "folder-outline",
        },
        {
            name: "profile",
            displayName: "Profile",
            icon: "settings-outline",
        },
    ];  

    return (
        <>
            <Tabs
                screenOptions={{
                    tabBarStyle: {
                        backgroundColor: "#fff",
                        borderTopWidth: 0,
                        elevation: 0,
                        height: 110,
                        paddingTop: 30,
                        zIndex: 40
                    },
                    headerShown: false,
                }}
            >
                {tabs.map(tab => (
                    <Tabs.Screen
                        key={tab.name}
                        name={tab.name}
                        options={{
                            tabBarLabel: () => null, 
                            tabBarIcon: ({ focused }) => (
                                <View style={{ alignItems: "center", justifyContent: "center", width: 100 }}>
                                    <Ionicons
                                        name={tab.icon as ComponentProps<typeof Ionicons>["name"]}
                                        size={28}
                                        color={focused ? "#e85a5a" : "#949598"}
                                    />
                                    <Text style={{
                                        fontSize: 12,
                                        marginTop: 4,
                                        fontWeight: "500",
                                        color: focused ? "#e85a5a" : "#949598"
                                    }}>
                                        {tab.displayName}
                                    </Text>
                                </View>
                            ),
                        }}
                    />
                ))}
            </Tabs>
        </>
    );
}