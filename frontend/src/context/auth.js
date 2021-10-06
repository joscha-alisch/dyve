import * as React from 'react'
import {useState,useEffect} from "react";
import {LoginPage} from "../pages/LoginPage";
import jwt_decode from "jwt-decode";
import {useCookies} from "react-cookie";
import useLocalStorage from "../hooks/useLocalStorage";
import axios from "axios";

const AuthContext = React.createContext()
const UserContext = React.createContext()

export const AuthProvider = ({children}) => {
    const [user, setUser] = useLocalStorage("user", null)

    const logout = async () => {
        axios.get("/auth/logout").then(invalidateAuth)
    }

    const invalidateAuth = () => {
        setUser(null)
    }

    axios.interceptors.response.use(response => {
        return response;
    }, error => {
        if (error.response.status === 401) {
            invalidateAuth()
        }
        return Promise.reject(error);
    });

    useEffect(() => {
        axios.get("/auth/user").then((res) => {
            setUser(res.data)
        })
    }, [])

    if (!user) {
        return <LoginPage />
    }

    return <AuthContext.Provider value={{logout, invalidateAuth}}>
            <UserContext.Provider value={user}>{children}</UserContext.Provider>
    </AuthContext.Provider>
}

export const useAuth = () => React.useContext(AuthContext)
export const useUser = () => React.useContext(UserContext)