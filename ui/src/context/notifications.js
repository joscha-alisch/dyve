import * as React from "react";
import {useEffect, useRef, useState} from "react";

const NotificationContext = React.createContext()

const TimedAlert = ({onClose, severity, message, timed = true, time, onTimeOut}) => {
    const [remaining, setRemaining] = useState(time)
    const [timer, setTimer] = useState()

    const resetTimer = () => {
        pauseTimer()
        setRemaining(time)
        const timer = setInterval(() => {
            setRemaining((previous) => previous - 100);
            if (remaining < 0) {
                console.log("yay")
                onTimeOut()
            }
        }, 100);
        setTimer(timer)
    }

    useEffect(() => {
        if (remaining < 0) {
            onTimeOut()
            pauseTimer()
        }
    }, [remaining])

    const pauseTimer = () => {
        clearInterval(timer)
        setRemaining(time)
    }

    useEffect(resetTimer, [])

    return <>
        {message}
    </>
}


export const NotificationProvider = ({children}) => {
    const [notifications, setNotifications] = useState([])

    const notify = (message, severity, manualClose = true) => {
        setNotifications([
            ...notifications,
            {message, severity, manualClose, open: true, key: Math.random()}
        ])
    }

    const dismiss = (index) => {
        const newList = [...notifications]
        newList[index].open = false
        setNotifications(newList)
    }

    const remove = (index) => {
        const newList = [...notifications]
        newList.splice(index, 1)
        setNotifications(newList)
    }

    return <NotificationContext.Provider value={{notify}}>
        {children}
        <ul>
            {notifications.map((notification, index) => <li key={notification.key}>
                    <div>
                        <TimedAlert onClose={() => dismiss(index)}
                                    severity={notification.severity}
                                    time={5000}
                                    message={notification.message} onTimeOut={() => dismiss(index)}/>
                    </div>
            </li>)}
        </ul>
    </NotificationContext.Provider>
}

export const useNotifications = () => React.useContext(NotificationContext)
