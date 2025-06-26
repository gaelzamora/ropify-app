import { useRouter } from "expo-router";
import { Image, StyleSheet, Text, TouchableOpacity, View } from "react-native";
import TypeWriter from 'react-native-typewriter'

export default function Index() {
    const router = useRouter()

    return (
        <View style={styles.container}>
            <View style={{ flex: 1, gap: 10, alignContent: "center", justifyContent: "center" }}>

                <View style={{ alignItems: "center" }}>
                    <Image 
                        source={require("@/assets/images/logo/logo_sinfondo.png")}
                        style={{ width: 100, height: 100 }}
                    />
                </View>

                <TypeWriter 
                    typing={1} 
                    style={styles.primaryText}
                >
                    Hello.
                </TypeWriter>
                <Text style={styles.secundaryText}>Lets Get Started!</Text>
        
                <TouchableOpacity 
                    style={[styles.button, { marginTop: 40 }]}
                    onPress={() => router.push("/login")}
                >
                    <Text style={styles.buttonText}>Sign in</Text>
                </TouchableOpacity>
                <Text style={{ textAlign: "center", fontWeight: "500" }}> or </Text>
                <TouchableOpacity 
                    style={styles.button}
                    onPress={() => router.push("/register")}
                >
                    <Text style={styles.buttonText}>Sign up</Text>
                </TouchableOpacity>
            </View>

        </View>
    )
}


const styles = StyleSheet.create({
    container: {
        backgroundColor: "#fff",
        flex: 1,
        alignItems: "center",
        justifyContent: "center"
    },
    containerText: {
        flex: 1,
        alignItems: "center",
        justifyContent: "center"
    },
    primaryText: {
        fontSize: 50,
        fontWeight: "500",
    },
    secundaryText: {
        fontSize: 30,
        fontWeight: "400"
    },
    button: {
        paddingVertical: 12,
        paddingHorizontal: 90,
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 4,
        elevation: 5,
        borderColor: "#d3d3d3",
        borderWidth: 1,
        borderRadius: 5,
    },
    buttonText: {
        fontWeight: "500"
    }
})