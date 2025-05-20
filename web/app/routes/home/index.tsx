import { useState } from "react";
import { HomePage } from "./component";

export default function Home() {
  // TODO implement
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
