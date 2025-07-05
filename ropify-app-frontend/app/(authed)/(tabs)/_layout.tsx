import { Ionicons } from "@expo/vector-icons";
import { Tabs } from "expo-router";
import React, { ComponentProps } from "react";
import { View, Text } from "react-native";

export default function TabLayout() {

    const tabs = [
        {
            name: "(feed)",
            displayName: "Feed",
            icon: "flame",
        },
        {
            name: "(closet)",
            displayName: "Closet",
            icon: "pricetag",
        },
        {
            name: "(outfit)",
            displayName: "Outfit",
            icon: "shirt",
        },
        {
            name: "profile",
            displayName: "Profile",
            icon: "settings",
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
                            tabBarIcon: ({ focused }) => {
                                // Si el icono tiene versión outline, la usamos cuando no está enfocado
                                const iconName = focused ? tab.icon : `${tab.icon}-outline`;
                                return (
                                    <View style={{ alignItems: "center", justifyContent: "center", width: 100 }}>
                                        <Ionicons
                                            name={iconName as ComponentProps<typeof Ionicons>["name"]}
                                            size={28}
                                            color={focused ? "#000" : "#888"}
                                        />
                                        <Text style={{
                                            fontSize: 12,
                                            marginTop: 4,
                                            fontWeight: "500",
                                            color: focused ? "#000" : "#888"
                                        }}>
                                            {tab.displayName}
                                        </Text>
                                    </View>
                                );
                            },
                        }}
                    />
                ))}
            </Tabs>
        </>
    );
}