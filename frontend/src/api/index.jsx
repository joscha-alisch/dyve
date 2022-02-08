import axios from "axios";
import {useCallback, useEffect, useState} from "react";
import useWebSocket, {ReadyState} from "react-use-websocket";
import {isDev} from "../helpers/isdev";

const Base = "/api"
const Paths = {
    Apps: {
        Get: (id) => "/apps/" + id,
        Live: (id) => "/apps/" + id + "/live",
    },
    Groups: {
        List: "/groups"
    },
    Teams: {
        Get:    (id) => "/teams/" + id,
        Create: (id) => "/teams/" + id
    }
}

const get = (api, then = () => {}, error) => axios.get(Base + api).then((res) => then(res.data))

const post = (api, data, then = () => {}, error = (err) => {}) => axios
    .post(Base + api, data)
    .then((res) => then(res.data))
    .catch(error)

const put = (api, data, then = () => {}, error = (err) => {}) => axios
    .put(Base + api, data)
    .then((res) => then(res.data))
    .catch(error)

const useWs = (api, setData) => {
    let socketUrl = "wss://" + window.location.host + Base + api
    if (isDev()) {
        socketUrl = "ws://localhost:9001" + Base + api
    }
    const {sendMessage, lastMessage, readyState } = useWebSocket(socketUrl, {
        shouldReconnect: (closeEvent) => true,
    });

    useEffect(() => {
        if (lastMessage !== null) {
            setData(JSON.parse(lastMessage.data))
        }
    }, [lastMessage, setData]);

    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];

    return {
        status: connectionStatus,
        send: sendMessage
    }
}

export const ListGroups = (then = () => {}, error = (err) => {}) => get(Paths.Groups.List, then, error)

export const GetApp = (id, then = () => {}, error = (err) => {}) => get(Paths.Apps.Get(id), then, error)

export const GetTeam = (id, then = () => {}, error = (err) => {}) => get(Paths.Teams.Get(id), then, error)
export const CreateTeam = (data, then = () => {}, error = (err) => {}) => post(Paths.Teams.Create(data.slug), data, then, error)

export const UseAppWebsocket = (id, setData) => useWs(Paths.Apps.Live(id), setData)

const API = {
    Apps: {
        Get: GetApp,
        Live: UseAppWebsocket,
    },
    Groups: {
        List: ListGroups
    },
    Teams: {
        Get: GetTeam,
        Create: CreateTeam
    }
}

export default API