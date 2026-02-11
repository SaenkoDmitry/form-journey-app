import "../styles/index.css";
import "../styles/Button.css";

type Props = React.ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: "primary" | "danger" | "ghost";
};

export default function Button({
                                   variant = "ghost",
                                   ...props
                               }: Props) {
    return (
        <button
            className={`btn btn-${variant}`}
            {...props}
        />
    );
}
