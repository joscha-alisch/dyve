import * as React from 'react'
import {useEffect} from 'react'
import useLocalStorage from "../hooks/useLocalStorage";
import axios from "axios";
import {useNotifications} from "./notifications";

const AuthContext = React.createContext()
export const UserContext = React.createContext()

export const AuthProvider = ({children}) => {
    const [user, setUser] = useLocalStorage("user", null)
    const {notify} = useNotifications()

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
            notify("Logged out", "error")
        }
        return Promise.reject(error);
    });

    useEffect(() => {
        axios.get("/auth/user").then((res) => {
            setUser(res.data)
        })
    }, [])

    if (!user) {
        return <span className="w-56 bg-red-200">not logged in</span>
    }

    return <AuthContext.Provider value={{logout, invalidateAuth}}>
        <UserContext.Provider value={user}>{children}</UserContext.Provider>
    </AuthContext.Provider>
}

export const useAuth = () => React.useContext(AuthContext)
export const useUser = () => React.useContext(UserContext)