import cn from "@/lib/utils";
import { useEffect, useState } from "react";

const navItems = [
    {name: "Home", href: "#home"},
    {name: "About", href: "#about"},
    {name: "Skills", href: "#skills"},
    {name: "Projects", href: "#projects"},
    {name: "Contact", href: "#contact"},
    {name: "Resume", href: "#Home"}
]

export const Navbar = () => {

    const [isScrolled, setIsScrolled] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
           setIsScrolled(window.scrollY > 10) 
        }
        window.addEventListener("scroll", handleScroll);
        return () => {
            window.removeEventListener("scroll", handleScroll);
        }
    }, [])

    return (
        <nav 
            className = {cn(
                "fixed w-full z-40 transition-all duration-300",
                isScrolled? "py-3 bg-background/80 backdrop-blur-md": "py-5 "
            )} 
        >
            <div className="container flex items-center justify-between"> 
                <a>
                    <span>
                        <span className="text-glow">Vaishnav</span>Portfolio 
                    </span>
                </a>

            </div>
        </nav>
    );
}