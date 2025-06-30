import { ActivityIndicator, Alert, FlatList, Image, StyleSheet, Text, TouchableOpacity, View, Button } from "react-native";
import React, { useCallback, useState } from "react";
import { FontAwesome, Ionicons } from "@expo/vector-icons";
import { Garment } from "@/types/garment";
import { useAuth } from "@/context/AuthContext";
import { garmentService } from "@/services/garment";
import { useFocusEffect } from "expo-router";
import Modal from 'react-native-modal'
import * as ImagePicker from 'expo-image-picker'
import {Camera} from 'expo-camera'

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
    const { user } = useAuth()
    const [isLoading, setIsLoading] = useState(false)

    // Scanner
    const [analyzing, setAnalyzing] = useState(false)
    const [garmentData, setGarmentData] = useState<Garment>()

    const [isModalActive, setIsModalActive] = useState(false)

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

    const takePhotoAndAnalyze = async () => {
        const { status } = await Camera.requestCameraPermissionsAsync()

        if (status !== 'granted') {
            Alert.alert('Se requieren permisos para la camara')
            return
        }

        let result = await ImagePicker.launchCameraAsync({
            allowsEditing: true,
            aspect: [4, 3],
            quality: 1,
        })

        if (result.canceled) {
            return
        }

        try {
            setAnalyzing(true)
            const imageUri = result.assets[0].uri

            // Enviar al backend para analisis

            const response = await garmentService.analyzeGarmentImage(imageUri)

            if (response.data && response.status === 200) {
                const { colors, category } = response.data

                if (garmentData?.id && garmentData.user_id && user) {
                    setGarmentData({
                        id: garmentData.id,
                        user_id: user.id,
                        name: garmentData?.name ?? "",
                        brand: garmentData?.brand ?? "",
                        size: garmentData?.size ?? "",
                        barcode: garmentData?.barcode ?? "",
                        is_verified: garmentData?.is_verified ?? "",
                        category: category || garmentData?.category || "",
                        color: colors && colors.length > 0 ? colors[0].name : garmentData?.color || "",
                        image_url: imageUri, // Guardar URI local temporalmente
                    })
                }
            }
        } catch (error) {
            console.error('Error analyzing image: ', error)
            Alert.alert('Error', 'No se pudo analizar la imagen')
        } finally {
            setAnalyzing(false)
        }
    }

    useFocusEffect(useCallback(() => { fetchClothes(activeClosetOption.toLowerCase()) }, [activeClosetOption]))

    return (
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
                            >
                                <Image 
                                    source={{ uri: garment.image_url }}
                                    style={styles.garmentImage}
                                />    
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
                <TouchableOpacity 
                    style={styles.iconTouchable}
                    onPress={() => setIsModalActive(true)}    
                >
                    <FontAwesome name="plus" size={25} color={"white"}  />
                </TouchableOpacity>
            </View>
            
            <Modal
                isVisible={isModalActive}
                onBackdropPress={() => setIsModalActive(false)}
                animationIn={"slideInUp"}
                animationOut={"zoomOut"}
                swipeDirection={"down"}
                onSwipeComplete={() => setIsModalActive(false)}
                backdropOpacity={0.6}
                style={styles.modal}
            >
                <View style={styles.modelContent}>
                    <Text>Este es un modal con gestos y animaciones</Text>
                    <Button title="Cerrar" onPress={() => setIsModalActive(false)} />
                </View>
            </Modal>
        </View>
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
    },
    garmentImage: {
        width: '100%',
        height: '100%',
        borderRadius: 15,
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
    }
})