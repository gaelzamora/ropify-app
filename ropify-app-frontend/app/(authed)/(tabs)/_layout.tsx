import { Ionicons } from "@expo/vector-icons";
import { Tabs } from "expo-router";
import React, { ComponentProps } from "react";
import { View, StyleSheet } from "react-native";

export default function TabLayout() {

    const tabs = [
        {
            name: "(home)",
            displayName: "Home",
            icon: "home-sharp",
        },
        {
            name: "(closet)",
            displayName: "Notes",
            icon: "document-text-sharp",
        },
        {
            name: "settings",
            displayName: "Settings",
            icon: "settings-outline",
        },
    ];  

    return (
        <>
            <Tabs
                screenOptions={{
                    tabBarStyle: {
                        backgroundColor: "transparent",
                        borderTopWidth: 0,
                        elevation: 0,
                        height: 120,
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
                                <View style={[
                                    styles.containerIcon,
                                    focused && styles.activeIcon
                                ]}>
                                    <Ionicons
                                        name={
                                             tab.name === "settings"
                                            ? (focused ? "settings-outline" : "settings-sharp")
                                            : tab.icon as ComponentProps<typeof Ionicons>["name"]   
                                        }
                                        size={28}
                                        color={focused ? "white" : "#949598"}
                                    />
                                </View>
                            ),
                        }}
                    />
                ))}
            </Tabs>
        </>
    );
}

const styles = StyleSheet.create({
    
    containerIcon: {
        height: 50,
        width: 50,
        alignItems: "center",
        justifyContent: "center",
        marginTop: 70,
    },
    activeIcon: {
        backgroundColor: "#4461ed",
        borderRadius: 40,
    }
});