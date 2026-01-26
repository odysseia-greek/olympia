import 'vuetify/styles';
import { createVuetify } from 'vuetify';
import '@mdi/font/css/materialdesignicons.css';

const vuetify = createVuetify({
    theme: {
        themes: {
            light: {
                colors: {
                    primary: '#1c61d1',
                    secondary: '#1cd18c',
                    triadic: '#1CBCD1',
                    papyrus: '#fdf6e3',
                    secondaryPapyrus: '#fefcf5',
                    accent: '#8c9eff',
                    text: '#fdf6e3',
                    textAccent: '#fefcf5',
                    error: '#d1311c',
                    footer: '#758b9a',
                    textprimary: "#394B6b",
                    textaccented: "#394B6b",
                    textsecondary: "#64789E",
                    background: '#e0f7fa',
                },
            },
        },
    },
});

export default vuetify;
