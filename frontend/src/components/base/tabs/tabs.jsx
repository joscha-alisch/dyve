import React, {useContext, useState} from "react"
import styles from "./tabs.module.sass"
import PropTypes from "prop-types"
import {TabContext as MuiTabContext, TabPanel as MuiTabPanel} from "@mui/lab";
import {Tab as MuiTab, Tabs as MuiTabs} from "@mui/material";

export const Tab = ({title, icon, children}) => children

const getMuiTabTitle = (child, index) => {
    if (child.type.name !== Tab.name) {
        console.warn("Only tab components allowed as children.");
        return ""
    }

    return <MuiTab label={child.props.title} value={index}/>
}

const getMuiTabContent = (child, index) => {
    return <MuiTabPanel sx={{padding: 0, paddingTop: "10px"}} value={"" + index} children={child.props.children}/>
}

const TabContext = React.createContext({})
export const useTabs = () => useContext(TabContext)

const Tabs = ({
                  className,
                  children,
                  renderHeader = (index, setIndex) => "",
                  renderFooter = ((index, setIndex) => "")
              }) => {
    let [index, setIndex] = useState(0)

    let providerValue = {
        current: index
    }

    return <div className={styles.Main + " " + className}>
        <TabContext.Provider value={providerValue}>
            <MuiTabContext value={"" + index}>
                <header className={styles.Header}>
                    <MuiTabs className={styles.Tabs} onChange={(e, v) => setIndex(v)} value={index}>
                        {React.Children.map(children, getMuiTabTitle)}
                    </MuiTabs>
                    <div className={styles.ExtraHeader}>
                        {renderHeader(index, setIndex)}
                    </div>
                </header>
                {React.Children.map(children, getMuiTabContent)}
                <footer className={styles.Footer}>
                    {renderFooter && renderFooter(index, setIndex)}
                </footer>
            </MuiTabContext>
        </TabContext.Provider>
    </div>
}

Tabs.propTypes = {
    className: PropTypes.string,
}

export default Tabs