import { router, useFocusEffect, useLocalSearchParams } from "expo-router";
import { View, Text, TouchableOpacity, FlatList, StyleSheet, RefreshControl, ActivityIndicator, Alert } from "react-native";
import { FontAwesome, Ionicons } from "@expo/vector-icons";
import SmartBackgroundRemoval from "@/components/SmartBackgroundRemoval";
import { useCallback, useState } from "react";
import { garmentService } from "@/services/garment";
import { Garment } from "@/types/garment";
import { Camera } from "expo-camera"
import * as ImagePicker from "expo-image-picker"
import Modal from "react-native-modal"
import { Image } from "expo-image"

const garmentCategories = [
    "all",
    "top",
    "bottom",
    "dress",
    "sneakers",
    "accesories",
    "backpack"
]

export default function ClosetDetailScreen() {
    const [activeClosetOption, setActiveClosetOption] = useState(garmentCategories[0])
    const [clothes, setClothes] = useState<Garment[]>([])
    const [isLoading, setIsLoading] = useState(false)
    const [isAnalyzing, setIsAnalyzing] = useState(false)
    const [refreshing, setRefreshing] = useState(false)

    const [elementsSelected, setElementsSelected] = useState<string[]>([])
    const [isDeleting, setIsDeleting] = useState(false)

    const [selectedGarment, setSelectedGarment] = useState<Garment | null>(null)
    const [isOpenGarment, setIsOpenGarment] = useState(false)
    const [closetName, setClosetName] = useState("...")

    const { id } = useLocalSearchParams()

    const onRefresh = useCallback(() => {
        setRefreshing(true)
        fetchClothes(activeClosetOption.toLowerCase())
            .finally(() => setRefreshing(false))
    }, [activeClosetOption])
    
    const fetchClothes = async (category: string) => {
        try {
            setIsLoading(true)
            const response = await garmentService.filterGarmentsFromCloset(id.toString(), 1, 20, category)
            setClothes(response.data.garments)
            setClosetName(response.data.closet_name)
        } catch (error) {
            Alert.alert("Error: ", String(error))
        } finally {
            setIsLoading(false)
        }
    }

    const fetchDeleteGarments = async (garments: string[]) => {
        try {
            await garmentService.deleteMultipleGarmentsFromCloset(garments)

            setElementsSelected([])
            setIsDeleting(false)

            await fetchClothes(activeClosetOption.toLowerCase())
        } catch (error) {
            Alert.alert("Error: ", String(error))
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
            const response = await garmentService.analyzeGarmentImage(id.toString(), imageUri)
            
            await fetchClothes(activeClosetOption.toLowerCase())
            if (response.status === 200) {
                Alert.alert('Success', 'Imagen cargada y procesada por IA')
            }
        } catch (error) {
            Alert.alert('Error', 'No se pudo analizar la imagen')
            console.error(error)
        } finally {
            setIsAnalyzing(false)
        }
    }

    useFocusEffect(useCallback(() => { fetchClothes(activeClosetOption.toLowerCase()) }, [activeClosetOption]))

    return (
        <>
            <View style={styles.closetContainer}>

                <View
                    style={{ position: "relative" }}
                >
                    <TouchableOpacity style={styles.backButton} onPress={() => {
                        router.back()}
                    }>
                        <FontAwesome name="arrow-left" size={25} color="#211" />
                    </TouchableOpacity>
                    
                    <Text style={styles.title}>{closetName}</Text>
                </View>

                
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
                                alignContent: "center",
                                justifyContent: 'flex-start',
                                width: "100%"
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
                                    <View style={{ flex: 1, marginTop: "50%", alignItems: "center", justifyContent: "center", padding: 40 }}>
                                        <ActivityIndicator size="large" color="#222" />
                                    </View>
                                ) : (
                                    <View style={{ flex: 1, marginTop: "50%", justifyContent: "center", alignItems: "center", padding: 40 }}>
                                        <FontAwesome name="tag" size={48} color="#7a7676" style={{ marginBottom: 10 }} />
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
                                    onPress={() => {
                                        if (isDeleting) {
                                            pushOnElementsSelected(garment.id)
                                        } else {
                                           setSelectedGarment(garment)
                                           setIsOpenGarment(true)
                                        }
                                    } 
                                }
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
                                                    color={"#222"}
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
                            onPress={() => fetchDeleteGarments(elementsSelected)}
                        >
                            <Ionicons 
                                name="trash-outline" 
                                size={25} 
                                color={"white"}  
                            />
                        </TouchableOpacity>
                    )}

                    {clothes.length > 0 && (
                        <TouchableOpacity 
                            style={[styles.iconTouchable, isDeleting ? {backgroundColor: "white"} : {}]}
                            onPress={() => {
                                if (isDeleting) setElementsSelected([])
                                setIsDeleting(!isDeleting)}
                            }   
                        >
                            <Ionicons 
                                name="ban-outline" 
                                size={25} 
                                color={isDeleting ? "#222" : "white"}  
                            />
                        </TouchableOpacity>
                    )}

                </View>
            

            </View>
                {isAnalyzing && (
                    <View style={styles.overlay}>
                        <ActivityIndicator size={"large"} color={"#fff"} />
                    </View>
                )}

            <Modal
                isVisible={isOpenGarment}
                onBackdropPress={() => {
                    setIsOpenGarment(false);
                    setTimeout(() => setSelectedGarment(null), 300);
                }}
                onSwipeComplete={() => {
                    setIsOpenGarment(false);
                    setTimeout(() => setSelectedGarment(null), 300);
                }}
                swipeDirection={['down']}
                backdropOpacity={0.7}
                animationIn="zoomIn"
                animationOut="zoomOut"
                animationInTiming={300}
                animationOutTiming={300}
                style={styles.modal}
            >
                <View style={styles.modalContent}>
                    {selectedGarment && (
                        <>
                            <TouchableOpacity 
                                style={styles.closeButton} 
                                onPress={() => {
                                    setIsOpenGarment(false);
                                    setTimeout(() => setSelectedGarment(null), 300);
                                }}
                            >
                                <Ionicons name="close" size={24} color="#fff" />
                            </TouchableOpacity>
                            
                            <View style={styles.garmentImageContainer}>
                                <Image
                                    source={selectedGarment.image_url}
                                    style={{ width: '100%', height: 300, borderRadius: 10 }}
                                    contentFit="contain"
                                    transition={300}
                                />
                            </View>
                        </>
                    )}
                </View>
            </Modal>
        </>
    )
}

const styles = StyleSheet.create({
    closetContainer: {
        flex: 1,
        paddingVertical: 60,
        paddingHorizontal: 20,
        position: 'relative'
    },
    title: {
        fontSize: 30, 
        fontWeight: "600",
        marginBottom: 25,
        textAlign: "center"
    },
    backButton: {
        position: "absolute",
        left: 0,
        top: 5
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
        height: "100%"
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
        borderColor: "#353333",
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
        backgroundColor: "#353333", 
        width: 50, 
        height: 50, 
        flex: 1, 
        alignItems: "center", 
        justifyContent: "center",
        borderRadius: 99
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
    modal: {
        margin: 0,
        justifyContent: 'center',
        alignItems: 'center',
    },
    modalContent: {
        backgroundColor: 'white',
        borderRadius: 15,
        padding: 20,
        width: '85%',
        maxHeight: '80%',
        alignItems: 'center',
    },
    closeButton: {
        position: 'absolute',
        top: -10,
        right: -10,
        backgroundColor: '#222',
        borderRadius: 15,
        width: 30,
        height: 30,
        justifyContent: 'center',
        alignItems: 'center',
        zIndex: 10,
    },
    garmentImageContainer: {
        width: '100%',
        height: 300,
        marginBottom: 20,
        borderRadius: 10,
        overflow: 'hidden',
    },
    garmentDetails: {
        width: '100%',
        padding: 10,
    },
    garmentName: {
        fontSize: 22,
        fontWeight: '700',
        marginBottom: 8,
        color: '#222',
    },
    garmentCategory: {
        fontSize: 16,
        color: '#666',
        marginBottom: 5,
    },
    garmentInfo: {
        fontSize: 14,
        color: '#777',
        marginBottom: 3,
    },
})