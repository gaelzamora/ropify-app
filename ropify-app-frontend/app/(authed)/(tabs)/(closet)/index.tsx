import { ActivityIndicator, Alert, FlatList, StyleSheet, Text, TouchableOpacity, View, RefreshControl, TextInput, Image, TouchableWithoutFeedback } from "react-native";
import React, { useCallback, useState } from "react";
import SmartBackgroundRemoval from "@/components/SmartBackgroundRemoval";
import { Closet } from "@/types/closet";
import { closetService } from "@/services/closet";
import { useAuth } from "@/context/AuthContext";
import { router, useFocusEffect } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import Modal from "react-native-modal"
import * as ImagePicker from "expo-image-picker"
import Toast from "react-native-toast-message";

export default function ClosetScreen() {
    // Closets
    const [closets, setClosets] = useState<Closet[]>([])
    const [isLoading, setIsLoading] = useState(false)
    const [isOpenAdd, setIsOpenAdd] = useState(false)

    const [closetName, setClosetName] = useState("");
    const [closetImage, setClosetImage] = useState<any>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const [isDeleteMode, setIsDeleteMode] = useState(false);
    const [closetToDelete, setClosetToDelete] = useState<string | null>(null);
    const DELETE_MODE_DELAY = 180;
    
    const { user } = useAuth()

    const [refreshing, setRefreshing] = useState(false)

    const handleLongPress = (id: string) => {
        setTimeout(() => {
            setIsDeleteMode(true);
            setClosetToDelete(id);
        }, DELETE_MODE_DELAY);
    };

    // Al salir del modo eliminación
    const handleExitDeleteMode = () => {
        setTimeout(() => {
            setIsDeleteMode(false);
            setClosetToDelete(null);
        }, DELETE_MODE_DELAY);
    };

    const onRefresh = useCallback(() => {
        setRefreshing(true)
        fetchClosets()
            .finally(() => setRefreshing(false))
    }, [])

    const fetchClosets = async () => {
        try {
            setIsLoading(true)
            const response = await closetService.getMany()
            setClosets(response.data)
        } catch (error) {
            Alert.alert("Error: ", String(error))
        } finally {
            setIsLoading(false)
        }
    }

    const deleteCloset = async (id: string) => {
        try {
            const response = await closetService.deleteOne(id);

            Toast.show({
                type: "success",
                text1: "Closet deleted",
                text2: response.statusText || "Closet deleted"
            })

            fetchClosets()
        } catch (error) {
            const err = error as any

            Toast.show({
                type: "error",
                text1: "Error",
                text2: err.response.data.message || "Oops, something went wrong, Please try again later"
            })
        }

    }

    const pickImage = async () => {
        let result = await ImagePicker.launchImageLibraryAsync({
            mediaTypes: ImagePicker.MediaTypeOptions.Images,
            allowsEditing: true,
            quality: 0.7,
        });
        if (!result.canceled && result.assets && result.assets.length > 0) {
            const asset = result.assets[0];
            setClosetImage({
                uri: asset.uri,
                name: asset.fileName || "image.jpg",
                type: asset.type || "image/jpeg"
            });
        }
    };

    // Handler para crear closet
    const handleCreateCloset = async () => {
        if (!closetName || !closetImage) {
            Alert.alert("Error", "Please provide a name and image.");
            return;
        }
        setIsSubmitting(true);
        try {
            const response = await closetService.createOne(closetName, closetImage);
            setIsOpenAdd(false);
            setClosetName("");
            setClosetImage(null);

            Toast.show({
                type: "success",
                text1: "Success",
                text2: response.message || "Closet created"
            })

            fetchClosets();
        } catch (err: any) {
            Toast.show({
                type: "error",
                text1: "Error",
                text2: err.response.data.message || "Oops, something went wrong, Please try again later"
            })
        } finally {
            setIsSubmitting(false);
        }
    };

    function onGoToClosetPage(id: string) {
        router.push(`/(authed)/(tabs)/(closet)/closet/${id}`)
    }

    useFocusEffect(useCallback(() => { fetchClosets() }, []))
    
    return (
        <View style={styles.closetContainer}>
            <View
                style={{ marginTop: 30 }}
            >
                {user?.firstName && user.lastName ? (
                    <Text
                        style={{ fontSize: 14, color: "#888", fontWeight: "400" }}
                    >Hi: {user?.firstName + ' ' + user?.lastName}</Text>
                ) : (
                    <Text
                        style={{ fontSize: 14, color: "#888", fontWeight: "400" }}
                    >Hi: {user?.email}</Text>
                )}
                <Text style={styles.textMain}>Closet&apos;s</Text>
            </View>

            <TouchableWithoutFeedback
                disabled={!isDeleteMode}
                onPress={() => {
                    handleExitDeleteMode()
                }}
            >
                <View style={{ flex: 1 }}>
                    <FlatList 
                        data={closets}
                        keyExtractor={(item) => item.id}
                        contentContainerStyle={{
                            alignContent: "center",
                            justifyContent: "center",
                            width: "100%",
                            gap: 10,
                            marginTop: 40
                        }}
                        ListEmptyComponent={
                            isLoading ? (
                                <View style={{ flex: 1, marginTop: "50%", alignItems: "center", justifyContent: "center", padding: 40 }}>
                                    <ActivityIndicator size={"large"} color={"#222"} />
                                </View>
                            ) : (
                                <View
                                    style={styles.isEmpty}
                                >
                                    <Ionicons 
                                        name="cube-outline"
                                        color={"#888"}
                                        size={66}
                                    />
                                    <Text style={{ textAlign: "center", fontSize: 24, fontWeight: "bold", color: "#888" }}>Is Empty here</Text>
                                </View>
                            )
                        }
                        refreshControl={
                            <RefreshControl 
                                refreshing={refreshing}  
                                onRefresh={onRefresh}
                                colors={["#222"]}
                                tintColor={"#222"}                    
                            />
                        }
                        renderItem={({ item: closet }) => (
                            <TouchableOpacity
                                style={styles.itemContainer}
                                onPress={() => onGoToClosetPage(closet.id)}
                                onLongPress={() => {
                                    handleLongPress(closet.id)
                                }}
                                activeOpacity={0.8}
                            >
                                <SmartBackgroundRemoval imageUri={closet.image_url} />
                                {isDeleteMode && closetToDelete === closet.id && (
                                    <View style={styles.deleteOverlay}>
                                        <TouchableOpacity
                                            style={styles.trashButton}
                                            onPress={async () => {
                                                deleteCloset(closet.id);
                                                setIsDeleteMode(false);
                                                setClosetToDelete(null);
                                                fetchClosets();
                                            }}
                                        >
                                            <Ionicons name="trash" size={32} color="#222" />
                                        </TouchableOpacity>
                                    </View>
                                )}
                            </TouchableOpacity>
                        )}
                    />
                </View>
            </TouchableWithoutFeedback>
            

            <TouchableOpacity
                style={styles.iconAddCloset}
                onPress={() => setIsOpenAdd(true)}
            >
                <Ionicons   
                    name={"add-outline"}
                    size={28}
                    color={"white"}
                />

            </TouchableOpacity>

            <Modal
                isVisible={isOpenAdd}
                onBackdropPress={() => setIsOpenAdd(false)}
                onSwipeComplete={() => setIsOpenAdd(false)}
                swipeDirection={['down']}
                backdropOpacity={0.7}
                animationIn="zoomIn"
                animationOut="zoomOut"
                animationInTiming={300}
                animationOutTiming={300}
                style={styles.modal}
            >
                <View style={{ width: 300, backgroundColor: "#fff", borderRadius: 12, padding: 20 }}>
                    <Text style={{ fontSize: 20, fontWeight: "bold", marginBottom: 10 }}>New Closet</Text>
                    <TextInput
                        placeholder="Name"
                        placeholderTextColor={"#888"}
                        value={closetName}
                        onChangeText={setClosetName}
                        style={{ borderWidth: 1, borderColor: "#ccc", borderRadius: 8, padding: 10, marginBottom: 15 }}
                    />
                    <TouchableOpacity
                        style={{
                            backgroundColor: "#222",
                            padding: 12,
                            borderRadius: 8,
                            alignItems: "center",
                            marginBottom: 10,
                        }}
                        onPress={pickImage}
                    >
                        <Text style={{ color: "#fff", fontWeight: "bold" }}>
                            {closetImage ? "Change Image" : "Pick Image"}
                        </Text>
                    </TouchableOpacity>
                    {closetImage && (
                        <View style={{ marginTop: "50%", alignItems: "center", marginBottom: 10 }}>
                            <Text style={{ fontSize: 12, marginBottom: 5 }}>{closetImage.name}</Text>
                            <Image 
                                source={ closetImage.uri } 
                                style={{ width: 80, height: 80, borderRadius: 8 }} 
                            />
                        </View>
                    )}
                    <TouchableOpacity
                        style={{
                            backgroundColor: "#222",
                            padding: 12,
                            borderRadius: 8,
                            alignItems: "center",
                            marginTop: 10
                        }}
                        onPress={handleCreateCloset}
                        disabled={isSubmitting}
                    >
                        <Text style={{ color: "#fff", fontWeight: "bold" }}>
                            {isSubmitting ? "Creating..." : "Create"}
                        </Text>
                    </TouchableOpacity>
                </View>
            </Modal>
            
        </View>            
    )
}

const styles = StyleSheet.create({
    closetContainer: {
        flex: 1,
        paddingVertical: 20,
        paddingHorizontal: 50,
        position: 'relative',
        backgroundColor: "white"
    },
    textMain: {
        fontSize: 40,
        fontWeight: "700"
    },
    itemContainer: {
        width: "100%",
        height: "auto",
    },
    iconAddCloset: {
        position: "absolute",
        bottom: 0,
        alignSelf: "center",
        backgroundColor: "#353333", 
        width: 50, 
        height: 50,
        justifyContent: "center",
        alignItems: "center",
        borderRadius: 99
    },
    modal: {
        margin: 0,
        justifyContent: 'center',
        alignItems: 'center',
    },
    deleteOverlay: {
        ...StyleSheet.absoluteFillObject,
        backgroundColor: "rgba(0,0,0,0.4)",
        justifyContent: "center",
        alignItems: "center",
        zIndex: 2,
        borderRadius: 25
    },
    trashButton: {
        backgroundColor: "#fff",
        padding: 18,
        borderRadius: 40,
        alignItems: "center",
        justifyContent: "center",
    },
    isEmpty: {
        flex: 1,
        justifyContent: "center",
        alignItems: "center",
        marginTop: "60%"
    }
})