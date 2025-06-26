import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios"
import { Platform } from "react-native"
import AsyncStorage from "@react-native-async-storage/async-storage"
import Constants from "expo-constants"

const url = Platform.OS === "android" ? Constants.expoConfig?.extra?.BASE_URL : Constants.expoConfig?.extra?.BASE_URL

const Api: AxiosInstance = axios.create({baseURL: url + "/api"})

Api.interceptors.request.use(async config => {
    const token = await AsyncStorage.getItem("token")

    if (token) config.headers.set("Authorization", `Bearer ${token}`)

    return config
})

Api.interceptors.response.use(
    async (res: AxiosResponse) => res.data,
    async (err: AxiosError) => Promise.reject(err)
)

export { Api }