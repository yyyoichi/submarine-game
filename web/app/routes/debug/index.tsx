import { useState } from "react";
import { HomePage } from "../home/component";

export default function Home() {
  const [isLoading, setIsLoading] = useState(false);
  const props: React.ComponentProps<typeof HomePage> = {
    StartButton: {
      disabled: isLoading,
      onClick: () => {
        setIsLoading(true);
        setTimeout(() => {
          setIsLoading(false);
        }, 1500);
      },
    },
  };
  return <HomePage {...props} />;
}
