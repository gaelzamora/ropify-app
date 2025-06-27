import { FlatList, StyleSheet, Text, TouchableOpacity, View } from "react-native";
import React, { useState } from "react";

const garmentCategories = [
    "all",
    "top",
    "bottom",
    "dress",
    "sneakers",
    "accesories",
    "backpack"
]

export default function ClosetScreen() {
    const [activeClosetOption, setActiveClosetOption] = useState(garmentCategories[0])

    return (
        <View style={styles.closetContainer}>
            <Text style={{ fontSize: 30, fontWeight: "600" }}>Closets</Text>
            
            <FlatList 
                data={garmentCategories}
                horizontal
                showsHorizontalScrollIndicator={false}
                keyExtractor={(item) => item}
                contentContainerStyle={{ paddingVertical: 10, gap: 10, marginTop: 15 }}
                renderItem={({ item }) => (
                    <TouchableOpacity
                        onPress={() => setActiveClosetOption(item)}
                        style={[styles.itemContainer, activeClosetOption === item && styles.itemActive]}
                    >
                        <Text 
                            style={[{ color: activeClosetOption === item ? "#" : "#777"}, styles.itemText]}
                        >
                            {item}
                        </Text>
                    </TouchableOpacity>
                )}
            />

        </View>
    )
}

const styles = StyleSheet.create({
    closetContainer: {
        paddingVertical: 30,
        paddingHorizontal: 20,
        width: "100%"
    },
    itemText: {
        fontWeight: "600", 
        textTransform: "capitalize", 
        textAlign: "center"
    },
    itemContainer: {
        paddingHorizontal: 10,
        paddingVertical: 6,
        alignItems: "center",
        justifyContent: "center",
    },
    itemActive: {
        borderColor: "#ee1e1e",
        borderBottomWidth: 3
    }
})