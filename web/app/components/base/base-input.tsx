import clsx from "clsx";
import { forwardRef, useEffect, useState } from "react";

interface Props {
  className?: string;
  value: string;
  lineHeight?: number;
  placeholder?: string;
  autoFocus?: boolean;
  onChange?: (value: string) => void;
  onEnter?: (value: string) => void;
  onFocus?: (focus: boolean) => void;
}

export const BaseInput = forwardRef<HTMLTextAreaElement, Props>(
  (props, ref) => {
    const {
      className,
      value,
      autoFocus,
      lineHeight = 24,
      onChange,
      onEnter,
    } = props;
    const [rows, setRows] = useState(1);

    const handleChange: React.ChangeEventHandler<HTMLTextAreaElement> = (
      event
    ) => {
      const value = event.target.value;
      onChange?.(value);

      const textareaLineHeight = 24;

      let wrapRows = value.split("\n").length;
      if (wrapRows >= 4) {
        setRows(4);
        return;
      }

      event.target.rows = 1;
      // 计算是不是换行
      const currentRows = Math.min(
        4,
        ~~(event.target.scrollHeight / textareaLineHeight)
      );

      setRows(currentRows);
      event.target.rows = currentRows;
    };

    useEffect(() => {
      if (!value) {
        setRows(1);
      }
    }, [value]);

    return (
      <textarea
        rows={rows}
        ref={ref}
        autoFocus={autoFocus}
        autoCapitalize="off"
        autoComplete="off"
        autoCorrect="off"
        maxLength={3600}
        placeholder={props.placeholder}
        onFocus={() => props.onFocus?.(true)}
        onBlur={() => props.onFocus?.(false)}
        className={clsx(
          "resize-none flex-1 mr-2 outline-none bg-transparent",
          "text-black/80 placeholder:text-black/40 caret-black/80",
          "dark:text-white/90 dark:placeholder:text-white/60 dark:caret-white/90",
          className
        )}
        style={{ lineHeight: `${lineHeight}px` }}
        value={value}
        onChange={handleChange}
        onKeyDown={(e) => {
          // https://www.zhangxinxu.com/wordpress/2023/02/js-enter-submit-compositionupdate/
          // https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent/isComposing
          if (!e.shiftKey && e.key === "Enter" && !e.nativeEvent.isComposing) {
            e.preventDefault();
            onEnter?.((e.target as any).value);
          }
        }}
      />
    );
  }
);
