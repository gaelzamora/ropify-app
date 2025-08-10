import React, { useCallback, useEffect, useRef, useState } from "react"
import { View, StyleSheet, FlatList, TouchableOpacity, ActivityIndicator, Alert, Text, Animated } from "react-native"
import { FontAwesome, Ionicons } from "@expo/vector-icons";
import Modal from 'react-native-modal'
import { outfitService } from "@/services/outfit";
import SmartBackgroundRemoval from "@/components/SmartBackgroundRemoval";
import { Garment, GarmentOptimized } from "@/types/garment";
import { Closet } from "@/types/closet";
import { closetService } from "@/services/closet";
import { useFocusEffect } from "expo-router";
import Toast from "react-native-toast-message";

type ComponentProps = {
    isModalOutfitGeneratedActive: boolean
    setIsModalOutfitGeneratedActive: (state: boolean) => void
}

type props = {
    visible: boolean
    message: string
    type: "success" | any
}

export const InModalToast = ({ visible, message, type = "success" }: props) => {
  const opacity = useRef(new Animated.Value(0)).current;
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        if (visible) {
        setMounted(true);
        Animated.timing(opacity, {
            toValue: 1,
            duration: 300,
            useNativeDriver: true,
        }).start();
        } else {
        Animated.timing(opacity, {
            toValue: 0,
            duration: 300,
            useNativeDriver: true,
        }).start(() => setMounted(false));
        }
    }, [visible]);

    if (!mounted) return null;

    return (
        <Animated.View style={{
        position: "absolute",
        top: 20,
        left: 20,
        right: 20,
        backgroundColor: type === "success" ? "#4CAF50" : "#F44336",
        padding: 10,
        borderRadius: 8,
        zIndex: 9999,
        flexDirection: "row",
        alignItems: "center",
        opacity,
        }}>
        <Ionicons
            name={type === "success" ? "checkmark-circle" : "alert-circle"}
            size={24}
            color="white"
            style={{ marginRight: 8 }}
        />
        <Text style={{ color: "white", fontWeight: "500" }}>{message}</Text>
        </Animated.View>
    );
};

export default function GenerateOutfitComponent({ isModalOutfitGeneratedActive, setIsModalOutfitGeneratedActive }: ComponentProps) {
    const [isLoadingRandomOutfit, setIsLoadingRandomOutfit] = useState(false)
    const [outfitRandom, setOutfitRandom] = useState<GarmentOptimized[]>([])
    const [closetSelected, setClosetSelected] = useState<string | undefined>(undefined)
    const [closets, setClosets] = useState<Closet[]>([])
    const [isLoadingClosets, setIsLoadingClosets] = useState(false)
    const [viewMode, setViewMode] = useState<"closets" | "outfit">("closets");
    const [isLoadingButton, setIsLoadingButton] = useState(false)

    const [inModalToast, setInModalToast] = useState({
        visible: false,
        message: "",
        type: "success"
    })

    const showInModalToast = (message: string, type = "success") => {
        setInModalToast({ visible: true, message, type })
        setTimeout(() => setInModalToast({ visible: false, message: "", type }), 3000)
    }

    // Cuando seleccionas un closet:
    const handleSelectCloset = (closetId: string) => {
        setClosetSelected(closetId);
        setViewMode("outfit");
        setOutfitRandom([]);
    };

    // Para volver a la vista de closets:
    const handleBackToClosets = () => {
        setViewMode("closets");
        setClosetSelected(undefined);
        setOutfitRandom([]);
    };

    const fetchGenerateRandomOutfit = async (closet_id: string, save: boolean) => {
        try {
            setIsLoadingRandomOutfit(true)
            const response = await outfitService.generateRandomOutfit(closet_id, save)
            setOutfitRandom(response.data.outfit.garments)

        } catch (error) {
            Alert.alert("Error: ", String(error))
        } finally {
            setIsLoadingRandomOutfit(false)
        }
    }

    const fetchClosets = async () => {
            try {
                setIsLoadingClosets(true)
                const response = await closetService.getMany()
                setClosets(response.data)
            } catch (error) {
                Alert.alert("Error: ", String(error))
            } finally {
                setIsLoadingClosets(false)
            }
    }

    const handleSaveOutfit = async () => {
        setIsLoadingButton(true)
        try {
            const response = await outfitService.createOutfit(outfitRandom, closetSelected, "")
            showInModalToast(response.message || "Outfit saved successfully")
        } catch (err: any) {
            showInModalToast(err.response?.data?.message || "Something went wrong", "error")
        } finally {
            setIsLoadingButton(false)
        }
    }

    useFocusEffect(useCallback(() => { fetchClosets() }, []))

    return (
        <Modal
            isVisible={isModalOutfitGeneratedActive}
            onBackdropPress={() => setIsModalOutfitGeneratedActive(false)}
            onSwipeComplete={() => setIsModalOutfitGeneratedActive(false)}
            swipeDirection={['down']}
            backdropOpacity={0.7}
            animationIn="slideInUp"
            animationOut="slideOutDown"
            animationInTiming={300}
            animationOutTiming={300}
            style={styles.modal}
        >
            <View style={styles.modalContent}>
                <InModalToast 
                    visible={inModalToast.visible}
                    message={inModalToast.message}
                    type={inModalToast.type}
                />
                <View style={styles.dragIndicator} />
                <View style={styles.contentRandomOutfit}>
                    {viewMode === "closets" && (
                        <>
                            <Text style={{ fontSize: 26, fontWeight: "bold", marginBottom: 20, color: "#222" }}>Select a closet</Text>
                            <FlatList
                                data={closets}
                                keyExtractor={(item) => item.id.toString()}
                                numColumns={2}
                                columnWrapperStyle={{ gap: 5 }}
                                contentContainerStyle={{
                                    alignContent: "center",
                                    justifyContent: "flex-start",
                                    width: "100%"
                                }}
                                ListEmptyComponent={
                                    isLoadingClosets ? (
                                        <View style={{ alignItems: "center", justifyContent: "center", padding: 20, marginTop: "50%" }}>
                                            <ActivityIndicator size={"large"} color={"#222"} />
                                        </View>
                                    ) : (
                                        <View style={{ alignItems: "center", justifyContent: "center", padding: 20, marginTop: "50%" }}>
                                            <FontAwesome name="random" size={48} color={"#888"} style={{ marginBottom: 10 }} />
                                            <Text style={styles.secundaryText}>Is Empty</Text>
                                        </View>
                                    )
                                }
                                renderItem={({ item: closet }) => (
                                    <TouchableOpacity
                                        onPress={() => handleSelectCloset(closet.id)}
                                        style={styles.garmentContainer}
                                    >
                                        <SmartBackgroundRemoval imageUri={closet.image_url} />
                                    </TouchableOpacity>
                                )}
                            />
                        </>
                    )}

                    {viewMode === "outfit" && (
                        <>
                            <TouchableOpacity
                                style={{ position: "absolute", left: 10, top: 0, padding: 10, zIndex: 10 }}
                                onPress={handleBackToClosets}
                            >
                                <Ionicons name="arrow-back" size={28} color="#222" />
                            </TouchableOpacity>
                            <FlatList
                                data={outfitRandom}
                                keyExtractor={(item) => item.id.toString()}
                                numColumns={3}
                                contentContainerStyle={{
                                    alignContent: "center",
                                    justifyContent: "flex-start",
                                    marginTop: 40
                                }}
                                ListEmptyComponent={
                                    isLoadingRandomOutfit ? (
                                        <View style={{ alignItems: "center", justifyContent: "center", padding: 20, marginTop: "50%" }}>
                                            <ActivityIndicator size={"large"} color={"#222"} />
                                        </View>
                                    ) : (
                                        <View style={{ alignItems: "center", justifyContent: "center", padding: 20, marginTop: "50%" }}>
                                            <FontAwesome name="random" size={48} color={"#888"} style={{ marginBottom: 10 }} />
                                            <Text style={styles.secundaryText}>Generate an outfit</Text>
                                        </View>
                                    )
                                }
                                renderItem={({ item: garment }) => (
                                    <TouchableOpacity style={[styles.garmentContainer, {flexBasis: "30%", maxWidth: "30%"}]}>
                                        <SmartBackgroundRemoval imageUri={garment.image_url} />
                                    </TouchableOpacity>
                                )}
                            />

                            <View
                                style={{ position: "absolute", bottom: 40, alignSelf: "center" }}
                            >
                                {closetSelected && outfitRandom.length > 0 && (
                                    <TouchableOpacity
                                        style={[styles.button, {bottom: 10}]}
                                        onPress={handleSaveOutfit}
                                    >
                                        {isLoadingButton ? (
                                            <ActivityIndicator color="#fff" size="small" />
                                        ) : (
                                            <Text style={{ color: "white", fontSize: 18, textAlign: "center" }}>Save Outfit</Text>
                                        )}
                                    </TouchableOpacity>
                                )}
                                {closetSelected && (
                                    <TouchableOpacity
                                        style={styles.button}
                                        onPress={() => fetchGenerateRandomOutfit(closetSelected, false)}
                                    >
                                        <Text style={{ color: "white", fontSize: 18 }}>Generate Outfit</Text>
                                    </TouchableOpacity>
                                )}
                            </View>
                        </>
                    )}
                </View>
            </View>
        </Modal>
    );
}


const styles = StyleSheet.create({
    modal: {
        margin: 0,
        justifyContent: 'flex-end',
        position: "relative"
    },
    modalContent: { 
        backgroundColor: 'white',
        borderTopLeftRadius: 15,
        borderTopRightRadius: 15,
        width: '100%',
        height: '80%',
        paddingTop: 15,
        position: "relative"
    },
    dragIndicator: {
        width: 60,
        height: 5,
        backgroundColor: '#ccc',
        borderRadius: 3,
        marginBottom: 10,
        alignSelf: 'center',
    },
    contentRandomOutfit: {
        flex: 1,
        paddingHorizontal: 14,
        position: "relative"
    },
    button: {
        paddingHorizontal: 100,
        paddingVertical: 16,
        backgroundColor: "#222",
        borderRadius: 10,
        bottom: 5,
        flexDirection: "row",
        justifyContent: "center",
        alignItems: "center"
    },
    garmentContainer: {
        aspectRatio: 1,
        margin: 5,
        backgroundColor: 'transparent', 
        alignItems: "center",
        justifyContent: "center",
        overflow: 'hidden',
        borderRadius: 15,
        position: "relative",
        flex: 1
    },
    secundaryText: {
        fontSize: 14, 
        color: "#888", 
        fontWeight: "600"
    }
})