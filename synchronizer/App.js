import * as React from 'react';
import { ScrollView, Text, TextInput, View, Image, TouchableHighlight, TouchableOpacity, StyleSheet } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import * as SecureStore from 'expo-secure-store';
import * as MediaLibrary from 'expo-media-library'
import axios from 'axios';
import { Button, Wrap, Box, Icon, IconComponentProvider } from "@react-native-material/core";
import jwtDecode from 'jwt-decode';
import mime from 'mime-types';
import MaterialCommunityIcons from "@expo/vector-icons/MaterialCommunityIcons";
import moment from 'moment';
import DateTimePicker from '@react-native-community/datetimepicker';
import pLimit from 'p-limit';
import Checkbox from 'expo-checkbox';

const STATUS_READY = 'ready'
const STATUS_INPROGRESS = 'inprogress'
const STATUS_SUCCESS = 'success'
const STATUS_FAILED = 'failed'
const STATUS_END = 'end'
const STATUS_SLEEPING = 'sleeping'
const STATUS_LOADING = 'loading'
const STATUS_DISPLAYING = 'displaying'
const STATUS_UNKNOWN = 'unknown'
const UPLOAD_MEDIAS = 'upload_medias'
const UPLOAD_FOLDER = 'upload_folder'

const API_HOST = ''


const AuthContext = React.createContext();

function SplashScreen() {
  return (
    <View>
      <Text>Loading...</Text>
    </View>
  );
}

function itemsReducer(state, action) {
  switch(action.type) {
    case 'add':
      return [...state, action.item]
    case 'remove':
      const items = []
  
      for (let it of state) {
        if (it.id === action.id) {
          it.toSync = !it.toSync
        }

        items.push(it)
      }
      
      return items
    case 'delete':
      return state.filter(({id}) => id !== action.id)
    
    case 'setSyncStatus':
      const items2 = []
      
      for (let it of state) {
        if (it.id === action.id) {
          it.syncStatus = action.status
        }

        items2.push(it)
      }

      return items2
    
    case 'reset':
      return []
    default:
      return state
  }
}

function HomeScreen() {
  const [lastSyncDate, setLastSyncDate] = React.useState()
  const [lastSyncDateTemp, setLastSyncDateTemp] = React.useState()
  const [syncStatus, setSyncStatus] = React.useState()
  const [items, setItems] = React.useReducer(itemsReducer, [])
  const [mode, setMode] = React.useState('date');
  const [show, setShow] = React.useState(false);
  const [loadFilesStatus, setLoadFilesStatus] = React.useState(STATUS_SLEEPING)
  const [globalSyncStatus, setGlobalSyncStatus] = React.useState(STATUS_SLEEPING)
  const [globalSyncStatusText, setGlobalSyncStatusText] = React.useState('')
  const [syncVideo, setSyncVideo] = React.useState(false);


  const getLastDate = async () => {
    let date = await SecureStore.getItemAsync("lastSyncDate")

    if (!date) {
      date = new Date(1672761025000).toISOString()
    }

    setLastSyncDate(date)
  }

  const listMedias = async () => {
    setLoadFilesStatus(STATUS_LOADING)
    try {
      await MediaLibrary.requestPermissionsAsync()

      const mediaTypes = ['photo']

      if (syncVideo) {
        mediaTypes.push('video')
      }

      let medias = await MediaLibrary.getAssetsAsync({
        mediaType: mediaTypes,
        first: 100,
        createdAfter: new Date(lastSyncDate).getTime()
      })

      setLoadFilesStatus(STATUS_DISPLAYING)

      let higherCreationTime = 0

      for (let {id, uri, filename, mediaType, creationTime} of medias.assets) {
        setItems({ type: 'add', item: {id, uri, filename, toSync: true, mediaType, syncStatus: "unknown", creationTime}} )

        if (creationTime > higherCreationTime) {
          higherCreationTime = creationTime
        }
      }

      if (medias.assets.length > 0) {
        setLastSyncDateTemp(new Date(higherCreationTime).toISOString())
      }

      setSyncStatus(STATUS_READY)
      setLoadFilesStatus(STATUS_SLEEPING)
    } catch (e) {
      alert(e)
    }
  }

  getLastDate()

  const mediaToSyncCount = () => {
    let v = items.reduce((ac, {toSync}) => {
      if (toSync) {
        return ac + 1
      }

      return ac
    }, 0)

    return <Text style={{color: 'white'}}>Synchroniser {v} medias</Text>
  }

  const endSyncText = () => {
    return <Text style={{color: 'red'}}>Terminer synchronisation</Text>
  }

  const startSyncText = () => {
    switch (loadFilesStatus) {
      case STATUS_LOADING:
        return <Text style={{color: 'white', textTransform: 'uppercase'}}>Chargement des medias...</Text>
      case STATUS_DISPLAYING:
        return <Text style={{color: 'white', textTransform: 'uppercase'}}>Affichage des medias</Text>
      default:
        return <Text style={{color: 'white', textTransform: 'uppercase'}}>Démarrer synchronisation</Text>
    }
  }

  const backSync = () => {
    setSyncStatus(STATUS_UNKNOWN)
    setItems({type: 'reset'})
  }

  const endSync = async () => {
    setSyncStatus(STATUS_END)
    setItems({type: 'reset'})
    setLastSyncDate(lastSyncDateTemp)


    try {
      await SecureStore.setItemAsync("lastSyncDate", lastSyncDateTemp)
    } catch(e) {
      console.log(e)
    }
  }

  const onChangeDate = async (event, selectedDate) => {
    const currentDate = new Date(selectedDate)
    setShow(false)
    setLastSyncDate(currentDate.toISOString())
    
    await SecureStore.setItemAsync("lastSyncDate", currentDate.toISOString())
    getLastDate()
  }

  const showMode = (currentMode) => {
    setShow(true);
    setMode(currentMode);
  }

  const showDatepicker = () => {
    showMode('date');
  }

  const synchronize = async () => {
    setGlobalSyncStatusText('')
    setGlobalSyncStatus(STATUS_INPROGRESS)
    const mediasConfig = items
    .filter(({toSync}) => toSync)
    .map(item => {
      return {
        name: item.filename,
        kind: item.mediaType.toUpperCase(),
        key: formatFileName(item.filename),
        uploadStatus: undefined,
        uri: item.uri,
      }
    })
    

    const querySignedUri = `
      query($medias: [GetIngestMediaInput!]!) {
        medias: ingest (input: {medias: $medias}) {
          key
          signedUri
        }
      }
    `

    const variablesSignedUri = {
      medias: mediasConfig.map(({key, kind}) => { return {key, kind} })
    }

    const authToken = await SecureStore.getItemAsync('authToken')

    try {
      const resp = await axios.post(
        `${API_HOST}/v3/graphql`,
        {
          query: querySignedUri,
          variables: variablesSignedUri
        },
        {
          headers:  { 'Authorization': 'Bearer ' + authToken }
        }
      )

      const promises = []
      const success = []
      const failed = []

      // Create a new array with the medias to upload

      const limit = pLimit(5)
      for (let media of resp.data.data.medias) {
          const conf = mediasConfig.find(({ key }) => key === media.key)

          let picture = await fetch(conf.uri)
          picture = await picture.blob()

          const mimeType = mime.lookup(media.key)

          let item = items.find(({filename}) => {
            return formatFileName(filename) === media.key
          })
          setItems({type: 'setSyncStatus', id: item.id, 'status': 'inprogress'})
          
          const req = fetch(
            media.signedUri,
            {
              method: 'PUT',
              body: picture,
              headers: {
                'Content-Type': mimeType,
                'Cache-Control': 'max-age=43200'
              }
            },
          ).then((resp) => {
            success.push(media.key)

            let item = items.find(({filename}) => {
              return formatFileName(filename) === media.key
            })

            setItems({type: 'setSyncStatus', id: item.id, 'status': 'success'})

          }).catch((e) => {
            failed.push(media.key)

            let item = items.find(({filename}) => {
              return formatFileName(filename) === media.key
            })

            setItems({type: 'setSyncStatus', id: item.id, 'status': 'failed'})
          })

          promises.push(limit(() => req))
      }

      setGlobalSyncStatusText('Envoi des médias sur le serveur')


      // Wait for all the medias to be uploaded
      // Then, create the folder and ingest the medias
      Promise.all(promises).then(async () => {
        const itemsUploaded = items.filter(({filename}) => {
            return success.includes(formatFileName(filename))
        })


        const authToken = await SecureStore.getItemAsync('authToken')
        const tokenDecoded = jwtDecode(authToken)
        const author = tokenDecoded.name

        let folders = {}

        const monthNames = ["janvier", "fevrier", "mars", "avril", "mai", "juin", "juillet", "aout", "septembre", "octobre", "novembre", "decembre"];
        for (let {id, creationTime, filename} of itemsUploaded) {
          const mediaDate = new Date(creationTime)
          let folderName = `${monthNames[mediaDate.getMonth()]}_${mediaDate.getFullYear()}_${author}`
          
          if (folders[folderName] === undefined) {
            folders[folderName] = []
          }

          let kind = 'PHOTO'
          if (filename.includes('mp4')) {
            kind = 'VIDEO'
          }

          folders[folderName].push({key: formatFileName(filename), kind, id})
        }

        setGlobalSyncStatusText('Création des dossiers')
        for (let folderName in folders) {
          const queryIngest = `
            mutation($medias: [PutIngestMediaInput!]!) {
              ingest(input: {medias: $medias}) {
                key
                status
              }
            }
          `

          const mediasUploaded = folders[folderName].map(({key, kind}) => {
            return {key, kind, author, folder: folderName}
          })

       
          const variablesIngest = {
            medias: mediasUploaded
          }

          try {
            await axios.post(
              `${API_HOST}/v3/graphql`,
              {
                query: queryIngest,
                variables: variablesIngest
              },
              {
                headers:  { 'Authorization': 'Bearer ' + authToken }
              }
            )

            for (let {id} of folders[folderName]) {
              setItems({type: 'delete', id})
            }
          } catch(e) {
            for (let {id} of folders[folderName]) {
              setItems({type: 'setSyncStatus', id, 'status': 'failed'})
            }
          }
        }

        setGlobalSyncStatus(STATUS_READY)
      })

    } catch(e) {
      console.log(e)
      setGlobalSyncStatus(STATUS_READY)
    }
  }

  const date = new Date(lastSyncDate)
  const formattedDate = moment(date).format('DD/MM/YYYY à HH:mm')

  return (
    <View style={{ padding: 20 }}>
      {items.length == 0 && syncStatus !== STATUS_READY && 
        <View>
          <Text style={{ fontSize: 20 }}>Dernière synchronisation le {formattedDate}</Text>
          <Button style={{ marginTop: 50 }} disabled={loadFilesStatus !== STATUS_SLEEPING} title={startSyncText} onPress={listMedias}/>
          <View style={{flexDirection: "row", marginTop: 20}}>
            <Checkbox value={syncVideo} onValueChange={setSyncVideo} />
            <TouchableOpacity onPress={() => setSyncVideo(!syncVideo)}>
              <Text style={{marginLeft: 20}}>Synchroniser les vidéos</Text>
            </TouchableOpacity>
          </View>
          <Button style={{marginTop: 70}} color="orange" title="Changer date synchronisation" onPress={showDatepicker} />
          {show && (
            <DateTimePicker 
              testID="dateTimePicker"
              value={new Date(lastSyncDate)}
              mode="date"
              is24Hour={true}
              onChange={onChangeDate}
            />
          )}
        </View>
      }
      <View>
        {syncStatus === 'ready' &&
          <Text style={{color: 'black'}} onPress={backSync}>
            <Icon name='keyboard-backspace' /> Retour
          </Text>
        }
        {items.length > 0 &&
        <View style={{marginTop: 5}}>
          <Text style={{ fontSize: 20, marginTop: 10 }}>{items.length} médias à synchronizer</Text>
          <ScrollView style={{height: 400, marginTop: 30}}>
            <Wrap style={{alignContent: 'center', justifyContent: 'center'}}>
              {items.map(({uri, id, filename, toSync, syncStatus, mediaType}) => {
                return (<TouchableHighlight key={filename} underlayColor="white" onPress={() => setItems({type: 'remove', id})}>
                  <Box w={100} h={100}>
                    {
                      syncStatus === STATUS_FAILED && <Box w={100} h={100} style={{position: 'absolute', zIndex: 200, backgroundColor: 'rgba(0, 0, 0, 0.4)'}}><Icon name='close' size={50} color="red" style={{position: 'absolute', top: 25, left: 25, zIndex: 200}}/></Box>
                      || syncStatus === STATUS_INPROGRESS && <Box w={100} h={100} style={{position: 'absolute', zIndex: 200, backgroundColor: 'rgba(0, 0, 0, 0.4)'}}><Icon name='cloud-upload' size={50} color="white" style={{position: 'absolute', top: 25, left: 25, zIndex: 200}}/></Box>
                      || syncStatus === STATUS_SUCCESS && <Box w={100} h={100} style={{position: 'absolute', zIndex: 200, backgroundColor: 'rgba(0, 0, 0, 0.4)'}}><Icon name='cloud-upload' size={50} color="green" style={{position: 'absolute', top: 25, left: 25, zIndex: 200}}/></Box>
                    }
                    {
                      <Box w={24} h={24} style={{position: 'absolute', bottom: 0, right: 0, zIndex: 200, backgroundColor: 'rgba(0, 0, 0, 0.4)'}}>
                        <Icon name={mediaType === 'photo' ? 'file-image' : 'file-video'} size={24} color="white" style={{zIndex: 200}}/>
                      </Box>
                    }
                    <Image source={{uri}} style={[styles.media, toSync ? styles.mediaEnabled : styles.mediaDisabled]}></Image>
                  </Box>
                </TouchableHighlight>)
              })}
            </Wrap>
          </ScrollView>
        </View>
        }
        {syncStatus === STATUS_READY && globalSyncStatus !== STATUS_INPROGRESS &&
          <View style={{marginTop: 20, flexDirection: 'column', justifyContent: 'center'}}>
              <Button disabled={globalSyncStatus === 'inprogress'} title={mediaToSyncCount} color="green" onPress={synchronize} style={{ marginRight: 20, marginLeft: 'auto', width: '100%' }} />
              <Button disabled={globalSyncStatus === 'inprogress'} title={endSyncText} style={{ backgroundColor: 'transparent', shadowColor: 'white', marginTop: 50 }} onPress={endSync} />
          </View>
        }
        {globalSyncStatus === STATUS_INPROGRESS && 
            <View style={{marginTop: 20, flexDirection: 'column', justifyContent: 'center'}}>
                <Text style={styles.uploadMessages}>{ globalSyncStatusText }</Text>
            </View>
        }
      </View>
    </View>
  );
}

function SignInScreen() {
  const [email, setEmail] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [signinInProgress, setSiginInProgress] = React.useState(false);
  const [signSuccess, setSignSuccess] = React.useState(null);

  const { signIn } = React.useContext(AuthContext);

  const sign = async () => {
    setSiginInProgress(true)
    setSignSuccess(null)
    const t = await signIn({ email, password })

    if (t !== 'ok') {
      setSignSuccess(false)
    }

    setSiginInProgress(false)
  }

  return (
    <View style={{padding: 20}}>
      {signSuccess === false && <Text style={{color: 'red'}}>Erreur de connexion</Text>}
      <TextInput
        placeholder="Email"
        value={email}
        onChangeText={setEmail}
        inputMode='email'
        style={styles.inputLogin}
      />
      <TextInput
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        style={styles.inputLogin}
      />
      <Button disabled={signinInProgress} title="Sign in" onPress={sign} />
    </View>
  );
}

function formatFileName(filename) {
  return filename.replace(/[^a-zA-Z0-9.]/g, '').normalize('NFD').replace(/[\u0300-\u036f]/g, '')
}

const Stack = createNativeStackNavigator();

export default function App({ navigation }) {
  const [state, dispatch] = React.useReducer(
    (prevState, action) => {
      switch (action.type) {
        case 'RESTORE_TOKEN':
          return {
            ...prevState,
            userToken: action.token,
            isLoading: false,
          };
        case 'SIGN_IN':
          return {
            ...prevState,
            isSignout: false,
            userToken: action.token,
          };
        case 'SIGN_OUT':
          return {
            ...prevState,
            isSignout: true,
            userToken: null,
          };
      }
    },
    {
      isLoading: true,
      isSignout: false,
      userToken: null,
    }
  );

  React.useEffect(() => {
    const bootstrapAsync = async () => {
      let authToken;

      try {
        authToken = await SecureStore.getItemAsync('authToken')
        
        const tokenDecoded = jwtDecode(authToken)

        const now = new Date()
        
        if (tokenDecoded.exp * 1000 < now.getTime()) {
          authContext.signOut()
        }

        dispatch({ type: 'RESTORE_TOKEN', token: authToken });

      } catch (e) {
        authContext.signOut()
      }
    };

    bootstrapAsync();
  }, []);

  const authContext = React.useMemo(
    () => ({
      signIn: async ({email, password}) => {
        const payload = `
        query {
          auth: auth(input: {email: "${email}", password: "${password}"}) {
            token
          }
        }
        `

        try {
          const resp = await axios.post(
            `${API_HOST}/v3/graphql`,
            {
              query: payload,
            }
          )
  
          if (resp.data.data.auth !== undefined) {
            await SecureStore.setItemAsync('authToken', resp.data.data.auth.token)
            dispatch({ type: 'SIGN_IN', token: resp.data.data.auth.token });
          }
        } catch (e) {
          console.log(e)

          return 'nok'
        }

        return 'ok'
      },
      signOut: async () => {
        await SecureStore.deleteItemAsync('authToken')

        dispatch({ type: 'SIGN_OUT' })
      },
    }),
    []
  );

  return (
    <IconComponentProvider IconComponent={MaterialCommunityIcons}>
      <AuthContext.Provider value={authContext}>
        <NavigationContainer>
          <Stack.Navigator>
            {state.isLoading && 1 == 2? (
              // We haven't finished checking for the token yet
              <Stack.Screen name="Splash" component={SplashScreen} />
            ) : state.userToken == null ? (
              // No token found, user isn't signed in
              <Stack.Screen
                name="Connexion"
                component={SignInScreen}
                options={{
                  title: 'Se connecter',
                  // When logging out, a pop animation feels intuitive
                  animationTypeForReplace: state.isSignout ? 'pop' : 'push',
                }}
              />
            ) : (
              // User is signed in
              <Stack.Screen name="AlbumSync" component={HomeScreen} />
            )}
          </Stack.Navigator>
        </NavigationContainer>
      </AuthContext.Provider>
    </IconComponentProvider>
  );
}

const styles = StyleSheet.create({
  media: {
    width: 100,
    height: 100,
    resizeMode: 'cover',
    overlayColor: 'red'
  },
  mediaDisabled: {
    opacity: 0.1
  },
  mediaEnabled: {
    opacity: 1
  },
  inputLogin: {
    borderStyle: 'solid',
    borderWidth: 1,
    borderColor: 'gray',
    borderRadius: 5,
    padding: 5,
    marginBottom: 20
  },
  uploadMessages: {
    textAlign: 'center',
  }
})
