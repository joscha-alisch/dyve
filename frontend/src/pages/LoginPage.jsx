import styles from "./LoginPage.module.sass"
import {useEffect, useState} from "react";
import {Spinner} from "../components/progress/Spinner";
import axios from "axios";
import {faKey, faLaptopCode} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'
import {faGithub} from "@fortawesome/free-brands-svg-icons";
import {isDev} from "../helpers/isdev";
import process from "process";


export const LoginPage = () => {
    let [providers, setProviders] = useState([])

    useEffect(() => {
        axios.get("/auth/list")
            .then(res => setProviders(res.data.map(provider => providerMap(provider))))
    }, [])

    return <div className={styles.Wrapper}>
        <img className={styles.Logo} alt="dyve logo" src="/img/logo.png" />
        <div className={styles.Box}>
            <h1>Login</h1>
            <ul className={styles.LoginProviders}>
                { providers.length === 0 ? <Spinner /> : ""}
                { providers.map((provider) => <li key={provider.name}>
                    <FontAwesomeIcon {...provider.icon} />
                    <a className={provider.className} href={provider.url}> {provider.name}</a>
                </li>) }
            </ul>
        </div>
    </div>
}

const providerUrl = (providerName) => {
    let host = isDev() ? "http://localhost:9001" : ""

    return host + "/auth/" + providerName + "/login?site=dyve&from=" + window.location.href
}
const providerMap = (providerName) => {
    switch (providerName) {
        case "dev":
            return {
                name: "Development",
                className: styles.ProviderGeneric,
                icon: { icon: faLaptopCode, transform: "grow-8" },
                url: providerUrl(providerName)
            }
        case "github":
            return {
                name: "GitHub",
                className: styles.ProviderGitHub,
                icon: { icon: faGithub, transform: "grow-10" },
                url: providerUrl(providerName)
            }
        default:
            return {
                name: providerName,
                className: styles.ProviderGeneric,
                icon: { icon: faKey, transform: "grow-6"},
                url: providerUrl(providerName)
            }
    }
}