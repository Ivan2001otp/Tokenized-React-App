import React from 'react'

export const TextFieldContainer = ({name,type,id,placeholder,value,onChanged})=>{
    return (
        <div className='mt-2'>

            <input

                id={id}
                name={name}
                type={type}
                autoComplete={id}
                required
                placeholder={placeholder}
                onChange={onChanged}
                value={value}
                className='block w-full rounded-md border-0 py-1.5 placeholder:px-2 text-gray-900 shadow-sm ring-1 ring-gray-300 placeholder:text-gray-400 focus:ring-insert focus:ring-indigo-600 sm:text-sm sm:leading-6'
            />
            

        </div>
    )
}