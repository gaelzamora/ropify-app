import { Ionicons } from "@expo/vector-icons";
import { Tabs } from "expo-router";
import React, { ComponentProps } from "react";
import { View, StyleSheet, Text } from "react-native";

export default function TabLayout() {

    const tabs = [
        {
            name: "(feed)",
            displayName: "Feed",
            icon: "flame",
        },
        {
            name: "(friends)",
            displayName: "Friends",
            icon: "people",
        },
        {
            name: "(closet)",
            displayName: "Closet",
            icon: "folder",
        },
        {
            name: "settings",
            displayName: "Settings",
            icon: "settings-sharp",
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
                        height: 100,
                        
                    },
                    headerShown: false,
                }}
            >
                {tabs.map(tab => (
                    <Tabs.Screen
                        key={tab.name}
                        name={tab.name}
                        options={{
                            tabBarLabel: ({ focused }) => (
                                <Text style={{ fontSize: 12, marginTop: 35, fontWeight: "500", color: focused ? "#e85a5a" : "#949598" }}>
                                    {tab.displayName}
                                </Text> 
                            ),
                            tabBarIcon: ({ focused }) => (
                                <View style={[
                                    styles.containerIcon,
                                ]}>
                                    <Ionicons
                                        name={
                                            tab.icon as ComponentProps<typeof Ionicons>["name"]   
                                        }
                                        size={28}
                                        color={focused ? "#e85a5a" : "#949598"}
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