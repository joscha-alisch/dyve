import * as React from "react";
import {useEffect, useRef, useState} from "react";
import styles from "./notifications.module.sass"
import {Alert, LinearProgress, Slide} from "@mui/material";

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
        <LinearProgress color={severity} className={styles.Timer} variant="determinate" value={(remaining / time) * 100} />
        <Alert onClose={onClose} onPointerEnter={pauseTimer} onPointerLeave={resetTimer}  severity={severity}>
            {message}
        </Alert>
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
        <ul className={styles.NotificationContainer}>
            {notifications.map((notification, index) => <li className={styles.Notification} key={notification.key}>
                <Slide direction={"left"} in={notification.open} onExited={() => remove(index)}>
                    <div>
                        <TimedAlert onClose={() => dismiss(index)}
                                    severity={notification.severity}
                                    time={5000}
                                    message={notification.message} onTimeOut={() => dismiss(index)}/>
                    </div>
                </Slide>
            </li>)}
        </ul>
    </NotificationContext.Provider>
}

export const useNotifications = () => React.useContext(NotificationContext)
