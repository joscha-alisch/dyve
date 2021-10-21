import {useAuth} from "../context/auth";
import {Redirect} from "react-router-dom";
import React from "react";

export const Logout = () => {
    const {logout} = useAuth()
    logout()

    return <Redirect to={"/"}/>
}