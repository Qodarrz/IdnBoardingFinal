import React from "react";

interface FormatedChatbotProps {
  text: string;
}

export const FormatedChatbot: React.FC<FormatedChatbotProps> = ({ text }) => {
  // Bold: **text**
  let formatted = text.replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>");

  // Underline: _text_
  formatted = formatted.replace(/_(.*?)_/g, "<u>$1</u>");

  // Bullet list: * text
  formatted = formatted.replace(/(?:\n)?\* (.*?)(?=\n|$)/g, "<li>$1</li>");

  // Wrap list items with <ul> if found
  if (formatted.includes("<li>")) {
    formatted = formatted.replace(/(<li>[\s\S]*<\/li>)/g, "<ul>$1</ul>");
  }

  // Line breaks
  formatted = formatted.replace(/\n/g, "<br />");

  return <div dangerouslySetInnerHTML={{ __html: formatted }} />;
};
