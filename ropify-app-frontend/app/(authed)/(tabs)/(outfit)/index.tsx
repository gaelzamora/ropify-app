import { ActivityIndicator, Alert, FlatList, StyleSheet, Text, TextInput, TouchableOpacity, View } from "react-native";
import React, { useState } from "react";
import GenerateOutfitComponent from "@/components/GenerateOutfitComponent";
import { Ionicons } from "@expo/vector-icons";


const tabs = [
    "My Outfits",
    "Recommended"
]

export default function OutfitScreen() {
    const [activeOptionSelected, setActiveOptionSelected] = useState(tabs[0])
    const [isModalOutfitGeneratedActive, setIsModalOutfitGeneratedActive] = useState(false)
    const [outfits, setOutfits] = useState<any>(null)
    const [searchQuery, setSearchQuery] = useState('')
    const [isLoading, setIsLoading] = useState(false)
    
    return (
        <>
            <View style={styles.outfitContainer}>
                <Text style={styles.mainText}>Outfits</Text>
            
                <View style={styles.contentArea}>
                    <View style={styles.optionSection}>
                        <View style={{ flexDirection: 'row', width: '100%' }}>
                            {tabs.map((item) => (
                                <TouchableOpacity
                                    key={item}
                                    onPress={() => setActiveOptionSelected(item)}
                                    style={[
                                        styles.itemContainer,
                                        activeOptionSelected === item && styles.itemActive,
                                        { flex: 1 }
                                    ]}
                                >
                                    <Text
                                        style={[
                                            { color: activeOptionSelected === item ? "#222" : "#777" }, 
                                            styles.itemText
                                        ]}
                                    >
                                        {item}
                                    </Text>
                                </TouchableOpacity>
                            ))}
                        </View>
                    </View>
                    
                    <View style={styles.searchContainer}>
                        <Ionicons 
                            name="search-outline" 
                            size={20} 
                            color="#999" 
                            style={styles.searchIcon} 
                        />
                        <TextInput 
                            style={styles.searchInput}
                            placeholder="Search"
                            placeholderTextColor="#999"
                            cursorColor={"#999"}
                            value={searchQuery}
                            onChangeText={setSearchQuery}
                        />
                    </View>
                    
                    <View style={styles.outfitSection}>
                        <FlatList
                            data={outfits}
                            keyExtractor={(item) => item}
                            numColumns={2}
                            contentContainerStyle={{
                                flex: 1,
                                justifyContent: 'flex-start',
                                gap: 10
                            }}
                            ListEmptyComponent={
                                isLoading ? (
                                    <View style={{ flex: 1, alignItems: "center", justifyContent: "center", padding: 40 }}>
                                        <ActivityIndicator size="large" color="#ee1e1e" />
                                    </View>
                                ) : (
                                    <View style={{ flex: 1, justifyContent: "center", alignItems: "center", padding: 40 }}>
                                        <Ionicons name="shirt" size={48} color="#7a7676" style={{ marginBottom: 10 }} />
                                        <Text style={{ fontSize: 20, color: "#7a7676", fontWeight: "700", textAlign: "center" }}>No clothes saved.</Text>
                                        <Text style={{ fontSize: 12, color: "#7a7676", textAlign: "center" }}>
                                            You haven&apos;t saved any clothes yet, so we don&apos;t have anything to show you! Go save some!.
                                        </Text>
                                    </View>
                                )
                            }
                            renderItem={({ item: any }) => (
                                <>
                                    <Text>Hello</Text>
                                </>
                            )}
                        />
                    </View>

                </View>
                    <View style={styles.containerButtons}>
                        <TouchableOpacity 
                            style={styles.buttonOutfit} 
                            onPress={() => setIsModalOutfitGeneratedActive(true)}
                        >
                            <Text style={styles.buttonText}>Generate Outfit</Text>
                        </TouchableOpacity>

                        <TouchableOpacity style={styles.buttonOutfit}>
                            <Text style={styles.buttonText}>Create Outfit</Text>
                        </TouchableOpacity>
                    </View>
                
                <GenerateOutfitComponent
                    isModalOutfitGeneratedActive={isModalOutfitGeneratedActive}
                    setIsModalOutfitGeneratedActive={setIsModalOutfitGeneratedActive}
                />

            </View>
        </>
    )
}

const styles = StyleSheet.create({
    outfitContainer: {
        flex: 1,
        paddingVertical: 50,
        paddingHorizontal: 20,
        position: 'relative'
    },
    mainText: {
        fontSize: 30,
        fontWeight: "600",
        marginBottom: 10
    },
    contentArea: {
        flex: 1,
        flexDirection: 'column',
    },
    optionSection: {
        height: 60,
        width: "100%"
    },
    itemText: {
        fontWeight: "600", 
        textTransform: "capitalize", 
        textAlign: "center"
    },
    itemContainer: {
        flex: 1, 
        paddingVertical: 6,
        height: 50,
        alignItems: "center",
        justifyContent: "center",
    },
    itemActive: {
        borderColor: "#353333",
        borderBottomWidth: 3
    },
    searchContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: '#ececec',
        borderRadius: 12,
        paddingVertical: 8,
        paddingHorizontal: 12,
        marginVertical: 12,
    },
    searchIcon: {
        marginRight: 8,
    },
    searchInput: {
        flex: 1,
        fontSize: 16,
        color: '#333',
        height: 30,
    },
    containerButtons: {
        flexDirection: "row",
        position: "absolute",
        bottom: 5,
        left: 18,
        right: 18,
        paddingVertical: 8,
        gap: 10,
    },
    buttonOutfit: {
        flex: 1,
        alignItems: "center",
        justifyContent: "center",
        paddingVertical: 15,
        backgroundColor: "#222",
        borderRadius: 8
    },
    buttonText: {
        color: "white",
        fontWeight: "600"
    },
    outfitSection: {
        flex: 1
    },
    
})