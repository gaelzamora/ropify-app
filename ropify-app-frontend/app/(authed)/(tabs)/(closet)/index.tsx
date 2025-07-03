import { ActivityIndicator, Alert, FlatList, StyleSheet, Text, TouchableOpacity, View, Button, RefreshControl } from "react-native";
import React, { useCallback, useState } from "react";
import { FontAwesome, Ionicons } from "@expo/vector-icons";
import { Garment } from "@/types/garment";
import { useAuth } from "@/context/AuthContext";
import { garmentService } from "@/services/garment";
import { useFocusEffect } from "expo-router";
import * as ImagePicker from 'expo-image-picker'
import {Camera} from 'expo-camera'
import SmartBackgroundRemoval from "@/components/SmartBackgroundRemoval";
import { Image } from "react-native";

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
    const [clothes, setClothes] = useState<Garment[]>([])
    const [isLoading, setIsLoading] = useState(false)
    const [isAnalyzing, setIsAnalyzing] = useState(false)
    const [refreshing, setRefreshing] = useState(false)

    const [elementsSelected, setElementsSelected] = useState<string[]>([])
    const [isDeleting, setIsDeleting] = useState(false)

    const { user } = useAuth()

    const onRefresh = useCallback(() => {
        setRefreshing(true)
        fetchClothes(activeClosetOption.toLowerCase())
            .finally(() => setRefreshing(false))
    }, [activeClosetOption])

    const fetchClothes = async (category: string) => {
        try {
            setIsLoading(true)
            const response = await garmentService.filterGarments(1, 12, user?.id, category)
            setClothes(response.data)
        } catch (error) {
            Alert.alert("Error: ", String(error))
        } finally {
            setIsLoading(false)
        }
    }

    const pushOnElementsSelected = (id: string) => {
        if (!id) return

        setElementsSelected((prevSelected: any) => {
            if (prevSelected.includes(id)) {
                return prevSelected.filter((item: string) => item !== id)
            } else {
                return [...prevSelected, id]
            }
        })
    }

    const takePhotoAndAnalyze = async () => {
        setIsAnalyzing(true)

        try {
            const { status } = await Camera.requestCameraPermissionsAsync()
            if (status !== 'granted') {
                setIsAnalyzing(false)
                Alert.alert('Se requieren permisos para la camara')
                return
            }
            const result = await ImagePicker.launchCameraAsync({
                allowsEditing: true,
                aspect: [4, 3],
                quality: 0.5,
            })
    
            if (result.canceled) {
                setIsAnalyzing(false)
                return
            }
            const imageUri = result.assets[0].uri
            const response = await garmentService.analyzeGarmentImage(imageUri)
            
            await fetchClothes(activeClosetOption.toLowerCase())
            if (response.status === 200) {
                Alert.alert('Success', 'Imagen cargada y procesada por IA')
            }
        } catch (error) {
            Alert.alert('Error', 'No se pudo analizar la imagen')
            console.log(error)
        } finally {
            setIsAnalyzing(false)
        }
    }

    useFocusEffect(useCallback(() => { fetchClothes(activeClosetOption.toLowerCase()) }, [activeClosetOption]))

    console.log(elementsSelected)

    return (
        <>
            <View style={styles.closetContainer}>
                
                <Text style={styles.title}>Closets</Text>
                
                <View style={styles.contentArea}>
                    <View style={styles.categorySection}>
                        <FlatList 
                            data={garmentCategories}
                            horizontal
                            showsHorizontalScrollIndicator={false}
                            keyExtractor={(item) => item}
                            contentContainerStyle={{ 
                                gap: 10
                            }}
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

                    <View style={styles.garmentSection}>
                        <FlatList 
                            data={clothes}
                            keyExtractor={(item) => item.id.toString()}
                            numColumns={3}
                            contentContainerStyle={{
                                flex: 1,
                                justifyContent: 'flex-start',
                                gap: 5,
                            }}
                            refreshControl={
                                <RefreshControl 
                                    refreshing={refreshing}
                                    onRefresh={onRefresh}
                                    colors={["#ee1e1e"]}
                                    tintColor={"#ee1e1e"}   
                                />
                            }
                            ListEmptyComponent={
                                isLoading ? (
                                    <View style={{ flex: 1, alignItems: "center", justifyContent: "center", padding: 40 }}>
                                        <ActivityIndicator size="large" color="#ee1e1e" />
                                    </View>
                                ) : (
                                    <View style={{ flex: 1, justifyContent: "center", alignItems: "center", padding: 40 }}>
                                        <FontAwesome name="save" size={48} color="#7a7676" style={{ marginBottom: 10 }} />
                                        <Text style={{ fontSize: 20, color: "#7a7676", fontWeight: "700", textAlign: "center" }}>No clothes saved.</Text>
                                        <Text style={{ fontSize: 12, color: "#7a7676", textAlign: "center" }}>
                                            You haven&apos;t saved any clothes yet, so we don&apos;t have anything to show you! Go save some!.
                                        </Text>
                                    </View>
                                )
                            }
                            renderItem={({ item: garment }) => (
                                <TouchableOpacity
                                    style={styles.garmentContainer}
                                    onPress={() => isDeleting ? pushOnElementsSelected(garment.id) : null}
                                >
                                    <SmartBackgroundRemoval
                                        imageUri={garment.image_url}
                                        boundingPoly={garment.boundingPoly}
                                    />
                                    {isDeleting && (
                                        <View style={styles.overlay}>
                                            {elementsSelected.includes(garment.id) && (
                                                <Ionicons 
                                                    name="checkmark-circle" 
                                                    size={25} 
                                                    color={"#ee1e1e"}
                                                    style={{position: "absolute", left: 2, top: 2}}
                                                />
                                            )}
                                        </View>
                                    )}
                                </TouchableOpacity>
                            )}
                        />
                    </View>
                </View>
                <View 
                    style={{ 
                        position: "absolute",
                        left: 15,
                        bottom: 15,
                    }}
                >
                    <TouchableOpacity 
                        style={styles.iconTouchable}
                        onPress={takePhotoAndAnalyze}
                    >
                        <Ionicons
                            name={"scan"}
                            size={28}
                            color={"white"}
                        />
                    </TouchableOpacity>
                </View>

                <View 
                    style={{ 
                        position: "absolute",
                        right: 15,
                        bottom: 15,
                    }}
                >
                    {elementsSelected.length > 0 && (
                        <TouchableOpacity 
                            style={[styles.iconTouchable, {bottom: 3}]}
                        >
                            <Ionicons 
                                name="trash-outline" 
                                size={25} 
                                color={"white"}  
                            />
                        </TouchableOpacity>
                    )}

                    <TouchableOpacity 
                        style={styles.iconTouchable}
                        onPress={() => setIsDeleting(!isDeleting)}
                    >
                        <Ionicons 
                            name="ban-outline" 
                            size={25} 
                            color={"white"}  
                        />
                    </TouchableOpacity>
                </View>
            

            </View>
                {isAnalyzing && (
                    <View style={styles.overlay}>
                        <ActivityIndicator size={"large"} color={"#fff"} />
                    </View>
                )}
        </>
    )
}

const styles = StyleSheet.create({
    closetContainer: {
        flex: 1,
        paddingVertical: 50,
        paddingHorizontal: 20,
        position: 'relative'
    },
    title: {
        fontSize: 30, 
        fontWeight: "600",
        marginBottom: 10
    },
    contentArea: {
        flex: 1,
        flexDirection: 'column',
    },
    categorySection: {
        height: 60, 
    },
    garmentSection: {
        flex: 1, 
    },
    itemText: {
        fontWeight: "600", 
        textTransform: "capitalize", 
        textAlign: "center"
    },
    itemContainer: {
        paddingHorizontal: 15,
        paddingVertical: 6,
        height: 50,
    },
    itemActive: {
        borderColor: "#ee1e1e",
        borderBottomWidth: 3
    },
    garmentContainer: {
        width: '30%',
        aspectRatio: 1,
        margin: 5,
        backgroundColor: 'transparent', 
        alignItems: "center",
        justifyContent: "center",
        overflow: 'hidden',
        borderRadius: 15,
        position: "relative"
    },
    garmentImage: {
        width: '100%',
        height: '100%',
        resizeMode: "cover",
        // Sombra para iOS:
        shadowColor: "#000",
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.2,
        shadowRadius: 4,
        // Sombra para Android:
        elevation: 5,
    },
    iconTouchable: {
        backgroundColor: "#e75959", 
        width: 50, 
        height: 50, 
        flex: 1, 
        alignItems: "center", 
        justifyContent: "center",
        borderRadius: 99
    },
    modal: {
        justifyContent: 'flex-end',
        margin: 0
    },
    modelContent: {
        backgroundColor:"white",
        padding: 20,
        borderTopLeftRadius: 20,
        borderTopRightRadius: 20,
        height: "80%"
    },
    processedBadge: {
        position: 'absolute',
        bottom: 5,
        right: 5,
        backgroundColor: 'rgba(238, 30, 30, 0.8)',
        borderRadius: 10,
        paddingHorizontal: 6,
        paddingVertical: 2,
    },
    processedText: {
        color: 'white',
        fontSize: 8,
        fontWeight: 'bold',
    },
    overlay: {
        position: "absolute",
        width: '100%',
        height: '100%',
        backgroundColor: 'rgba(0,0,0,0.3)',
        alignItems: 'center',
        justifyContent: 'center',
  },
})